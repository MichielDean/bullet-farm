# Context

## Item: ci-vlpsw

**Title:** Add docs cataracta between QA and delivery
**Status:** in_progress
**Priority:** 2

### Description

Add a documentation writer cataracta that runs between qa and delivery. Reviews the diff, finds documentation gaps, updates or creates docs in the repo.

Changes required:

1. cataractae_definitions: add docs_writer role. Instructions: read CONTEXT.md; run git diff main...HEAD; find all .md files; for each changed area (CLI, config, pipeline, architecture) check if docs exist and are accurate; if no user-visible changes pass with 'No documentation updates required'; otherwise update outdated sections, add missing docs, commit with '<id>: docs: update documentation for changes'; signal pass with file list or recirculate with specific ambiguity.

2. cataractae list: insert docs step between qa and delivery. Fields: name=docs, type=agent, identity=docs_writer, model=sonnet, context=full_codebase, skills=[cistern-droplet-state], timeout_minutes=20, skip_for=[1], on_pass=delivery, on_fail=implement, on_recirculate=implement, on_escalate=human

3. Routing: qa on_pass changes delivery -> docs. Trivial complexity skip_cataractae gains docs.

4. Mirror: copy updated aqueduct/feature.yaml to cmd/ct/assets/aqueduct/feature.yaml - both identical.

5. Tests: update any tests referencing cataracta list or pipeline step names.

## Current Step: implement

- **Type:** agent
- **Role:** implementer
- **Context:** full_codebase

## Recent Step Notes

### From: manual

Implemented docs_writer cataracta between qa and delivery. Committed: 9aab594. Tests: all 9 packages pass. Verified: docs_writer definition in feature.yaml (2 occurrences), docs step with skip_for=[1]/on_pass=delivery, security-review on_pass updated to docs, trivial skip_cataractae gains docs, mirrors identical (diff = empty), cataractae/docs_writer/CLAUDE.md created, scheduler_test.go updated (complexityWorkflow+docs step, TestAdvanceSkipped_StandardSkipsQA expects docs not delivery, TestComplexity_CriticalHumanGateBeforeMerge+HumanGateSetsCurrentCataracta start at docs not qa).

### From: manual

Phase 2: aqueduct/aqueduct.yaml and cmd/ct/assets/aqueduct/aqueduct.yaml — docs step is missing on_escalate: human. The diff adds on_escalate: human only to the security-review step; the docs step inherits the old security-review context lines (on_pass: delivery, on_fail: implement, on_recirculate: implement) but not on_escalate. Requirement explicitly lists on_escalate=human for docs; scheduler_test.go line 229 includes OnEscalate: human in the workflow struct. Tests pass because the workflow is constructed directly in code, not parsed from YAML. At runtime, critical droplets (complexity=4) at the docs step will not escalate to human review before delivery. Fix: add on_escalate: human to the docs step in both YAML files.

### From: manual

Already implemented — no changes required. Verified: both aqueduct/aqueduct.yaml:447 and cmd/ct/assets/aqueduct/aqueduct.yaml:447 have on_escalate: human for the docs step. scheduler_test.go:1050 has OnEscalate: 'human'. All 9 packages pass.

### From: scheduler

Implement pass rejected: HEAD has not advanced since last review (commit: 00cd796d50e7457551093dc04850db43ea2cd705). No new commits were found. You must commit your changes before signaling pass.

<available_skills>
  <skill>
    <name>cistern-droplet-state</name>
    <description>Manage droplet state in the Cistern agentic pipeline using the `ct` CLI.</description>
    <location>.claude/skills/cistern-droplet-state/SKILL.md</location>
  </skill>
  <skill>
    <name>github-workflow</name>
    <description>---</description>
    <location>.claude/skills/github-workflow/SKILL.md</location>
  </skill>
</available_skills>

## Signaling Completion

When your work is done, signal your outcome using the `ct` CLI:

**Pass (work complete, move to next step):**
    ct droplet pass ci-vlpsw

**Recirculate (needs rework — send back upstream):**
    ct droplet recirculate ci-vlpsw
    ct droplet recirculate ci-vlpsw --to implement

**Block (genuinely blocked, cannot proceed):**
    ct droplet block ci-vlpsw

Add notes before signaling:
    ct droplet note ci-vlpsw "What you did / found"

The `ct` binary is on your PATH.
