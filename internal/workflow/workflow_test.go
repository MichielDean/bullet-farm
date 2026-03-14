package workflow

import (
	"path/filepath"
	"strings"
	"testing"
)

func testdataPath(name string) string {
	return filepath.Join("testdata", name)
}

func TestParseValidWorkflow(t *testing.T) {
	w, err := ParseWorkflow(testdataPath("valid_workflow.yaml"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Name != "feature" {
		t.Errorf("name = %q, want %q", w.Name, "feature")
	}
	if len(w.Steps) != 4 {
		t.Fatalf("got %d steps, want 4", len(w.Steps))
	}

	impl := w.Steps[0]
	if impl.Name != "implement" {
		t.Errorf("step[0].Name = %q, want %q", impl.Name, "implement")
	}
	if impl.Type != StepTypeAgent {
		t.Errorf("step[0].Type = %q, want %q", impl.Type, StepTypeAgent)
	}
	if impl.Role != "implementer" {
		t.Errorf("step[0].Role = %q, want %q", impl.Role, "implementer")
	}
	if impl.Model != "sonnet" {
		t.Errorf("step[0].Model = %q, want %q", impl.Model, "sonnet")
	}
	if impl.Context != ContextFullCodebase {
		t.Errorf("step[0].Context = %q, want %q", impl.Context, ContextFullCodebase)
	}
	if impl.MaxIterations != 3 {
		t.Errorf("step[0].MaxIterations = %d, want 3", impl.MaxIterations)
	}
	if impl.TimeoutMinutes != 30 {
		t.Errorf("step[0].TimeoutMinutes = %d, want 30", impl.TimeoutMinutes)
	}
	if impl.OnPass != "review" {
		t.Errorf("step[0].OnPass = %q, want %q", impl.OnPass, "review")
	}
	if impl.OnFail != "blocked" {
		t.Errorf("step[0].OnFail = %q, want %q", impl.OnFail, "blocked")
	}

	review := w.Steps[1]
	if review.OnRevision != "implement" {
		t.Errorf("step[1].OnRevision = %q, want %q", review.OnRevision, "implement")
	}
	if review.OnEscalate != "human" {
		t.Errorf("step[1].OnEscalate = %q, want %q", review.OnEscalate, "human")
	}

	merge := w.Steps[3]
	if merge.Type != StepTypeAutomated {
		t.Errorf("step[3].Type = %q, want %q", merge.Type, StepTypeAutomated)
	}
}

func TestCircularRouteError(t *testing.T) {
	_, err := ParseWorkflow(testdataPath("circular_route.yaml"))
	if err == nil {
		t.Fatal("expected circular route error, got nil")
	}
	if !strings.Contains(err.Error(), "circular route") {
		t.Errorf("error = %q, want it to contain 'circular route'", err)
	}
}

func TestMissingRefError(t *testing.T) {
	_, err := ParseWorkflow(testdataPath("missing_ref.yaml"))
	if err == nil {
		t.Fatal("expected missing ref error, got nil")
	}
	if !strings.Contains(err.Error(), "nonexistent-step") {
		t.Errorf("error = %q, want it to mention 'nonexistent-step'", err)
	}
	if !strings.Contains(err.Error(), "unknown step") {
		t.Errorf("error = %q, want it to contain 'unknown step'", err)
	}
}

func TestUnknownTypeError(t *testing.T) {
	_, err := ParseWorkflow(testdataPath("unknown_type.yaml"))
	if err == nil {
		t.Fatal("expected unknown type error, got nil")
	}
	if !strings.Contains(err.Error(), "unknown type") {
		t.Errorf("error = %q, want it to contain 'unknown type'", err)
	}
	if !strings.Contains(err.Error(), "magic") {
		t.Errorf("error = %q, want it to mention 'magic'", err)
	}
}

func TestMaxIterationsError(t *testing.T) {
	_, err := ParseWorkflow(testdataPath("bad_max_iterations.yaml"))
	if err == nil {
		t.Fatal("expected max_iterations error, got nil")
	}
	if !strings.Contains(err.Error(), "max_iterations") {
		t.Errorf("error = %q, want it to contain 'max_iterations'", err)
	}
}

func TestParseWorkflowBytes(t *testing.T) {
	yaml := `
name: simple
steps:
  - name: do-thing
    type: gate
    on_pass: done
`
	w, err := ParseWorkflowBytes([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Name != "simple" {
		t.Errorf("name = %q, want %q", w.Name, "simple")
	}
	if len(w.Steps) != 1 {
		t.Fatalf("got %d steps, want 1", len(w.Steps))
	}
	if w.Steps[0].Type != StepTypeGate {
		t.Errorf("type = %q, want %q", w.Steps[0].Type, StepTypeGate)
	}
}

func TestValidateEmptyName(t *testing.T) {
	w := &Workflow{Steps: []WorkflowStep{{Name: "x", Type: StepTypeAgent}}}
	err := Validate(w)
	if err == nil || !strings.Contains(err.Error(), "name is required") {
		t.Errorf("expected name required error, got %v", err)
	}
}

func TestValidateNoSteps(t *testing.T) {
	w := &Workflow{Name: "empty"}
	err := Validate(w)
	if err == nil || !strings.Contains(err.Error(), "no steps") {
		t.Errorf("expected no steps error, got %v", err)
	}
}

func TestValidateDuplicateStepName(t *testing.T) {
	w := &Workflow{
		Name: "dup",
		Steps: []WorkflowStep{
			{Name: "a", Type: StepTypeAgent, OnPass: "done"},
			{Name: "a", Type: StepTypeAgent, OnPass: "done"},
		},
	}
	err := Validate(w)
	if err == nil || !strings.Contains(err.Error(), "duplicate step name") {
		t.Errorf("expected duplicate step error, got %v", err)
	}
}

func TestTerminalRefsAreValid(t *testing.T) {
	// "done", "blocked", "human", "escalate" should be accepted as targets.
	yaml := `
name: terminals
steps:
  - name: s1
    type: agent
    on_pass: done
    on_fail: blocked
    on_revision: human
    on_escalate: escalate
`
	_, err := ParseWorkflowBytes([]byte(yaml))
	if err != nil {
		t.Fatalf("terminal refs should be valid, got: %v", err)
	}
}
