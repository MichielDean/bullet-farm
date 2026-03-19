# Context

## Item: ci-vwi46

**Title:** Add rate limiting to the delivery cataracta API endpoint
**Status:** in_progress
**Priority:** 1

### Description

Prevent abuse of the droplet ingestion endpoint. Apply per-IP and per-token limits with configurable thresholds in cistern.yaml.

## Current Step: implement

- **Type:** agent
- **Role:** implementer
- **Context:** full_codebase

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
    ct droplet pass ci-vwi46

**Recirculate (needs rework — send back upstream):**
    ct droplet recirculate ci-vwi46
    ct droplet recirculate ci-vwi46 --to implement

**Block (genuinely blocked, cannot proceed):**
    ct droplet block ci-vwi46

Add notes before signaling:
    ct droplet note ci-vwi46 "What you did / found"

The `ct` binary is on your PATH.
