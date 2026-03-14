package runner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/MichielDean/bullet-farm/internal/bd"
	"github.com/MichielDean/bullet-farm/internal/workflow"
)

func testRepoConfig() workflow.RepoConfig {
	return workflow.RepoConfig{
		Name:     "testrepo",
		URL:      "https://github.com/example/testrepo",
		Workers:  3,
		Names:    []string{"alice", "bob", "charlie"},
		BdPrefix: "tr-",
	}
}

func testWorkflow() *workflow.Workflow {
	return &workflow.Workflow{
		Name: "feature",
		Steps: []workflow.WorkflowStep{
			{Name: "implement", Type: "agent", Role: "implementer", Context: "full_codebase"},
			{Name: "review", Type: "agent", Role: "reviewer", Context: "diff_only"},
		},
	}
}

func TestNewRunner_NamedWorkers(t *testing.T) {
	cfg := Config{
		Repo:        testRepoConfig(),
		Workflow:    testWorkflow(),
		BdClient:    bd.NewClient("bd", ""),
		SandboxRoot: t.TempDir(),
	}

	r, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	workers := r.Workers()
	if len(workers) != 3 {
		t.Fatalf("expected 3 workers, got %d", len(workers))
	}

	want := []struct {
		name, session string
	}{
		{"alice", "testrepo-alice"},
		{"bob", "testrepo-bob"},
		{"charlie", "testrepo-charlie"},
	}

	for i, w := range want {
		if workers[i].Name != w.name {
			t.Errorf("worker %d: name = %q, want %q", i, workers[i].Name, w.name)
		}
		if workers[i].SessionID != w.session {
			t.Errorf("worker %d: session = %q, want %q", i, workers[i].SessionID, w.session)
		}
		if workers[i].Repo != "testrepo" {
			t.Errorf("worker %d: repo = %q, want %q", i, workers[i].Repo, "testrepo")
		}
	}
}

func TestNewRunner_NumberedWorkers(t *testing.T) {
	repo := testRepoConfig()
	repo.Names = nil // No namepool — use numbered names.

	cfg := Config{
		Repo:        repo,
		Workflow:    testWorkflow(),
		BdClient:    bd.NewClient("bd", ""),
		SandboxRoot: t.TempDir(),
	}

	r, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	workers := r.Workers()
	for i, w := range workers {
		expected := workerName(repo, i)
		if w.Name != expected {
			t.Errorf("worker %d: name = %q, want %q", i, w.Name, expected)
		}
	}
}

func TestNewRunner_NoWorkflow(t *testing.T) {
	cfg := Config{
		Repo:     testRepoConfig(),
		BdClient: bd.NewClient("bd", ""),
	}
	_, err := New(cfg)
	if err == nil {
		t.Fatal("expected error for nil workflow")
	}
}

func TestNewRunner_NoBdClient(t *testing.T) {
	cfg := Config{
		Repo:     testRepoConfig(),
		Workflow: testWorkflow(),
	}
	_, err := New(cfg)
	if err == nil {
		t.Fatal("expected error for nil bd client")
	}
}

func TestClaimRelease(t *testing.T) {
	cfg := Config{
		Repo:        testRepoConfig(),
		Workflow:    testWorkflow(),
		BdClient:    bd.NewClient("bd", ""),
		SandboxRoot: t.TempDir(),
	}

	r, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if r.IdleCount() != 3 {
		t.Fatalf("idle count = %d, want 3", r.IdleCount())
	}

	// Claim all workers.
	w1 := r.Claim()
	w2 := r.Claim()
	w3 := r.Claim()
	if w1 == nil || w2 == nil || w3 == nil {
		t.Fatal("expected 3 claims to succeed")
	}

	if r.IdleCount() != 0 {
		t.Fatalf("idle count = %d, want 0", r.IdleCount())
	}

	// No more workers available.
	w4 := r.Claim()
	if w4 != nil {
		t.Fatal("expected nil when all workers busy")
	}

	// Release one.
	r.Release(w2)
	if r.IdleCount() != 1 {
		t.Fatalf("idle count = %d, want 1", r.IdleCount())
	}

	// Claim again.
	w5 := r.Claim()
	if w5 == nil {
		t.Fatal("expected claim after release")
	}
	if w5.Name != w2.Name {
		t.Errorf("expected released worker %q, got %q", w2.Name, w5.Name)
	}
}

func TestStepByName(t *testing.T) {
	cfg := Config{
		Repo:        testRepoConfig(),
		Workflow:    testWorkflow(),
		BdClient:    bd.NewClient("bd", ""),
		SandboxRoot: t.TempDir(),
	}

	r, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	step := r.StepByName("implement")
	if step == nil {
		t.Fatal("expected step 'implement'")
	}
	if step.Role != "implementer" {
		t.Errorf("step role = %q, want %q", step.Role, "implementer")
	}

	if r.StepByName("nonexistent") != nil {
		t.Error("expected nil for nonexistent step")
	}
}

func TestOutcomeValidate(t *testing.T) {
	tests := []struct {
		result  string
		wantErr bool
	}{
		{"pass", false},
		{"fail", false},
		{"revision", false},
		{"escalate", false},
		{"", true},
		{"unknown", true},
	}

	for _, tt := range tests {
		o := &Outcome{Result: tt.result}
		err := o.Validate()
		if (err != nil) != tt.wantErr {
			t.Errorf("Validate(%q): err=%v, wantErr=%v", tt.result, err, tt.wantErr)
		}
	}
}

func TestOutcomeRouteField(t *testing.T) {
	tests := []struct {
		result string
		want   string
	}{
		{"pass", "on_pass"},
		{"fail", "on_fail"},
		{"revision", "on_revision"},
		{"escalate", "on_escalate"},
		{"unknown", ""},
	}

	for _, tt := range tests {
		o := &Outcome{Result: tt.result}
		if got := o.RouteField(); got != tt.want {
			t.Errorf("RouteField(%q) = %q, want %q", tt.result, got, tt.want)
		}
	}
}

func TestWriteContextFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "CONTEXT.md")

	bead := &bd.Bead{
		ID:          "bf-123",
		Title:       "Test bead",
		Status:      "in_progress",
		Priority:    1,
		Description: "Fix the thing",
	}

	step := &workflow.WorkflowStep{
		Name:    "implement",
		Type:    "agent",
		Role:    "implementer",
		Context: "full_codebase",
	}

	notes := []bd.StepNote{
		{FromStep: "review", Text: "Looks good but needs tests"},
	}

	err := writeContextFile(path, ContextParams{
		Level:      "full_codebase",
		SandboxDir: dir,
		Bead:       bead,
		Step:       step,
		Notes:      notes,
	})
	if err != nil {
		t.Fatalf("writeContextFile: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read CONTEXT.md: %v", err)
	}

	content := string(data)
	checks := []string{
		"# Context",
		"bf-123",
		"Test bead",
		"implementer",
		"From: review",
		"Looks good but needs tests",
		"outcome.json",
	}
	for _, want := range checks {
		if !contains(content, want) {
			t.Errorf("CONTEXT.md missing %q", want)
		}
	}
}

func TestPrepareContext_FullCodebase(t *testing.T) {
	dir := t.TempDir()
	bead := &bd.Bead{ID: "bf-1", Title: "Test", Status: "open", Priority: 1}
	step := &workflow.WorkflowStep{Name: "implement", Type: "agent", Context: "full_codebase"}

	ctxDir, cleanup, err := PrepareContext(ContextParams{
		Level:      workflow.ContextFullCodebase,
		SandboxDir: dir,
		Bead:       bead,
		Step:       step,
	})
	if err != nil {
		t.Fatalf("PrepareContext: %v", err)
	}
	defer cleanup()

	if ctxDir != dir {
		t.Errorf("ctxDir = %q, want sandbox dir %q", ctxDir, dir)
	}

	if _, err := os.Stat(filepath.Join(ctxDir, "CONTEXT.md")); err != nil {
		t.Error("expected CONTEXT.md in sandbox dir")
	}
}

func TestPrepareContext_SpecOnly(t *testing.T) {
	bead := &bd.Bead{
		ID:          "bf-2",
		Title:       "Spec test",
		Status:      "open",
		Priority:    2,
		Description: "Build the widget",
	}
	step := &workflow.WorkflowStep{Name: "plan", Type: "agent", Context: "spec_only"}

	ctxDir, cleanup, err := PrepareContext(ContextParams{
		Level:      workflow.ContextSpecOnly,
		SandboxDir: t.TempDir(),
		Bead:       bead,
		Step:       step,
	})
	if err != nil {
		t.Fatalf("PrepareContext: %v", err)
	}
	defer cleanup()

	// spec.md should exist.
	specData, err := os.ReadFile(filepath.Join(ctxDir, "spec.md"))
	if err != nil {
		t.Fatal("expected spec.md")
	}
	if !contains(string(specData), "Spec test") {
		t.Error("spec.md missing bead title")
	}

	// CONTEXT.md should exist.
	if _, err := os.Stat(filepath.Join(ctxDir, "CONTEXT.md")); err != nil {
		t.Error("expected CONTEXT.md")
	}
}

func TestOutcomeJSON(t *testing.T) {
	dir := t.TempDir()
	outcome := Outcome{
		Result: "pass",
		Notes:  "All tests passing",
		Annotations: map[string]string{
			"tests_added": "3",
		},
	}

	data, err := json.Marshal(outcome)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	path := filepath.Join(dir, "outcome.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	readData, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	var parsed Outcome
	if err := json.Unmarshal(readData, &parsed); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if parsed.Result != "pass" {
		t.Errorf("result = %q, want %q", parsed.Result, "pass")
	}
	if parsed.Annotations["tests_added"] != "3" {
		t.Error("missing annotation tests_added")
	}
}

func TestHandoffPrepend(t *testing.T) {
	dir := t.TempDir()

	// Write initial CONTEXT.md.
	ctxPath := filepath.Join(dir, "CONTEXT.md")
	if err := os.WriteFile(ctxPath, []byte("# Original Context\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Write handoff.md.
	handoffPath := filepath.Join(dir, "handoff.md")
	if err := os.WriteFile(handoffPath, []byte("## Progress\nDid step A, need step B.\n"), 0644); err != nil {
		t.Fatal(err)
	}

	sess := &Session{ID: "test-sess", WorkDir: dir}
	if err := sess.prependHandoffToContext(); err != nil {
		t.Fatalf("prependHandoffToContext: %v", err)
	}

	data, err := os.ReadFile(ctxPath)
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)
	if !contains(content, "Handoff from Previous Session") {
		t.Error("missing handoff header")
	}
	if !contains(content, "Did step A") {
		t.Error("missing handoff content")
	}
	if !contains(content, "# Original Context") {
		t.Error("missing original context")
	}
}

func TestWorkerSandboxPaths(t *testing.T) {
	sandboxRoot := "/tmp/test-sandboxes"
	repo := testRepoConfig()

	workers, err := initWorkers(repo, filepath.Join(sandboxRoot, repo.Name))
	if err != nil {
		t.Fatalf("initWorkers: %v", err)
	}

	expected := []string{
		filepath.Join(sandboxRoot, "testrepo", "alice"),
		filepath.Join(sandboxRoot, "testrepo", "bob"),
		filepath.Join(sandboxRoot, "testrepo", "charlie"),
	}

	for i, w := range workers {
		if w.SandboxDir != expected[i] {
			t.Errorf("worker %d: sandbox = %q, want %q", i, w.SandboxDir, expected[i])
		}
	}
}

func TestWorkerName_Fallback(t *testing.T) {
	repo := workflow.RepoConfig{
		Name:    "test",
		Workers: 2,
		Names:   []string{"alpha"},
	}

	if got := workerName(repo, 0); got != "alpha" {
		t.Errorf("workerName(0) = %q, want %q", got, "alpha")
	}
	if got := workerName(repo, 1); got != "worker-1" {
		t.Errorf("workerName(1) = %q, want %q", got, "worker-1")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
