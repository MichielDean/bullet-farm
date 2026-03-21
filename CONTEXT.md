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

## Current Step: implement

- **Type:** agent
- **Role:** implementer
- **Context:** full_codebase

## Recent Step Notes

### From: manual

Phase 2: dashboard_web_test.go readWSTextFrame case 127 — the 8-byte extended payload length is parsed as only the lower 32 bits (ext[4]<<24 | ext[5]<<16 | ext[6]<<8 | ext[7]); ext[0]–ext[3] (the high 32 bits) are silently discarded, violating RFC 6455 §5.2. Fix: length = int(ext[0])<<56 | int(ext[1])<<48 | int(ext[2])<<40 | int(ext[3])<<32 | int(ext[4])<<24 | int(ext[5])<<16 | int(ext[6])<<8 | int(ext[7]). Dead code in current tests (all payloads <126 bytes) but wrong by spec and would silently produce a bad length for any future large-frame test. No other issues found.

### From: manual

Fixed readWSTextFrame case 127 in dashboard_web_test.go: 8-byte extended payload length was only using low 32 bits (ext[4]–ext[7]); corrected to include all 8 bytes per RFC 6455 §5.2. All 9 packages pass.

### From: manual

Fixed readWSTextFrame case 127: 8-byte extended length now correctly combines all 8 bytes per RFC 6455 §5.2. All 9 packages pass.

### From: scheduler

Implement pass rejected: HEAD has not advanced since last review (commit: adf8afd900b608cb93db0d1d6b0998ddde882e2d). No new commits were found. You must commit your changes before signaling pass.

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
    ct droplet pass ci-0vm8f

**Recirculate (needs rework — send back upstream):**
    ct droplet recirculate ci-0vm8f
    ct droplet recirculate ci-0vm8f --to implement

**Block (genuinely blocked, cannot proceed):**
    ct droplet block ci-0vm8f

Add notes before signaling:
    ct droplet note ci-0vm8f "What you did / found"

The `ct` binary is on your PATH.
