# Context

## Item: ci-0vm8f

**Title:** Cataractae peek: read-only live observer for active aqueduct sessions
**Status:** in_progress
**Priority:** 2

### Description

Add ability to observe any active cataractae session in real-time without interacting with it. Requirements:
- GET /api/aqueducts/{name}/peek returns current tmux pane content as text
- WebSocket endpoint /ws/aqueducts/{name}/peek streams live pane output (poll tmux every 500ms, send diffs)
- Web UI: clicking an active aqueduct arch opens a peek panel/modal showing live session output
- Read-only: no keyboard input forwarded, no interaction possible, purely observational
- Shows last N lines of pane (configurable, default 100)
- Auto-scrolls to bottom, toggle to pin scroll position
- Clear label: 'Observing — read only'
- Falls back gracefully if aqueduct is idle or tmux session not found

## Current Step: simplify

- **Type:** agent
- **Role:** simplifier
- **Context:** full_codebase

## Recent Step Notes

### From: manual

Simplified: removed unused prevContent field from peekModel (dead code), renamed cap->capturer in peek HTTP/WS handlers to avoid shadowing Go builtin, replaced nested bounds checks with max()/min(). Tests: all 9 packages pass.

### From: manual

Phase 2: (1) cmd/ct/peek_tui.go Init() — double polling logic bug: tea.Batch(m.fetchCmd(), peekTickCmd()) starts two independent 500ms polling loops. fetchCmd() fires immediately, peekContentMsg handler schedules peekTickCmd() (loop A); the initial peekTickCmd() in Init() fires at t=500ms, its peekTickMsg handler calls fetchCmd(), which returns peekContentMsg, which schedules another peekTickCmd() (loop B). Both loops run indefinitely, polling tmux at ~2x the intended rate (~250ms effective interval). Fix: Init() should return only peekTickCmd() or only fetchCmd(), not both. (2) cmd/ct/dashboard_web.go — new HTTP/WS behavior has zero test coverage: lookupAqueductSession, GET /api/aqueducts/{name}/peek handler, WS /ws/aqueducts/{name}/peek handler, wsUpgrade, wsSendText, and wsAcceptKey are all untested. peek_tui_test.go covers only the TUI model.

### From: manual

Fixed double polling bug in peekModel.Init(): was calling tea.Batch(fetchCmd(), peekTickCmd()) which started two independent 500ms loops; changed to return only peekTickCmd() for a single tick→fetch→tick loop. Added 13 new tests in dashboard_web_test.go covering wsAcceptKey (RFC 6455 test vector), wsSendText small+medium frame encoding, lookupAqueductSession (empty/no-match/found), GET /api/aqueducts/{name}/peek (405, idle, active with mock capturer, session-gone, ?lines= param), and WS endpoint rejection (426 for non-WS, 400 for missing key). All 9 packages pass. Committed 9d4d96f.

### From: manual

Fixed double polling in Init() (single peekTickCmd() loop). Added 13 tests covering wsAcceptKey, wsSendText (small+medium), lookupAqueductSession, GET peek endpoint (405/idle/active-mock/gone-session/?lines), WS rejection (426/400). All packages pass. Committed 9d4d96f.

<available_skills>
  <skill>
    <name>cistern-droplet-state</name>
    <description>Manage droplet state in the Cistern agentic pipeline using the `ct` CLI.</description>
    <location>.claude/skills/cistern-droplet-state/SKILL.md</location>
  </skill>
  <skill>
    <name>code-simplifier</name>
    <description>code-simplifier</description>
    <location>.claude/skills/code-simplifier/SKILL.md</location>
  </skill>
</available_skills>

## Signaling Completion

When your work is done, signal your outcome using the `ct` CLI:

**Pass (work complete, move to next step):**
    ct droplet pass ci-0vm8f

**Recirculate (needs rework — send back upstream):**
    ct droplet recirculate ci-0vm8f
    ct droplet recirculate ci-0vm8f --to implement

**Block (genuinely blocked, cannot proceed):**
    ct droplet block ci-0vm8f

Add notes before signaling:
    ct droplet note ci-0vm8f "What you did / found"

The `ct` binary is on your PATH.
