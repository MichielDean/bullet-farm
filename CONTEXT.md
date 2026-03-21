# Context

## Item: ci-xilu9

**Title:** Bug: in-repo skills not available in non-cistern repo sandboxes
**Status:** in_progress
**Priority:** 2

### Description

When aqueducts run against ScaledTest or PortfolioWebsite sandboxes, skills defined as in-repo paths (e.g. path: skills/cistern-droplet-state/SKILL.md) fail to copy because those paths only exist in the cistern repo, not in ScaledTest or PortfolioWebsite.

Root cause: internal/cataractae/runner.go resolves skill.Path relative to w.SandboxDir. For ScaledTest sandboxes, w.SandboxDir is ~/.cistern/sandboxes/ScaledTest/julia — which doesn't contain skills/cistern-droplet-state/.

Fix options:
1. Treat skills with path: as external skills if the path doesn't exist in the sandbox, falling back to ~/.cistern/skills/<name>/SKILL.md
2. When ct repo add creates new sandboxes, copy the cistern repo's skills/ directory into each new sandbox
3. Change aqueduct.yaml to use external skill references (no path:) for universal skills like cistern-droplet-state, and install them to ~/.cistern/skills/ during ct init

Option 1 is the most robust — graceful fallback means it works regardless of how the sandbox was created.

Affected: every ScaledTest and PortfolioWebsite aqueduct run logs 'warning: copy skill cistern-droplet-state' and agents run without their state management skill.

## Current Step: simplify

- **Type:** agent
- **Role:** simplifier
- **Context:** full_codebase

## Recent Step Notes

### From: manual

REVISED APPROACH: deterministic, no agent path reliance.

The fix should be infrastructure-level, not prompt-level:

1. runner.go: when an in-repo skill (path: set) doesn't exist in the sandbox, fall back to ~/.cistern/skills/<name>/SKILL.md before giving up. This makes the copy deterministic — agent always gets the file.

2. ct init: seed all built-in skills (cistern-droplet-state, github-workflow, code-simplifier) to ~/.cistern/skills/ so the fallback location is always populated.

3. ct doctor: add check that every skill referenced in aqueduct.yaml exists in ~/.cistern/skills/. Flag missing skills as warnings with instructions to run ct skills install.

4. ct repo add: after cloning, copy ~/.cistern/skills/ into the new sandbox's .claude/skills/ directory so it's seeded at creation time.

The injected <location> in context.go should remain .claude/skills/<name>/SKILL.md — an agent-relative path that is guaranteed to exist because the runner always copies it there. No absolute paths, no agent guessing, no reliance on context.

### From: manual

FINAL APPROACH: use --add-dir to give Claude access to ~/.cistern/skills, inject absolute path in CONTEXT.md.

Two changes:

1. internal/cataractae/session.go line 55: add --add-dir ~/.cistern/skills to the claude command:
   claudeCmd := fmt.Sprintf('%s --dangerously-skip-permissions --add-dir ~/.cistern/skills %s-p ...', ...)

2. internal/cataractae/context.go line 258: change the injected <location> from the relative sandbox path to the absolute skills store path:
   <location>~/.cistern/skills/<name>/SKILL.md</location>

3. The copy step in runner.go (lines 223-248) can be removed entirely — no files need to be copied since Claude has direct read access to the source.

4. ct init: ensure ~/.cistern/skills/ is seeded with all built-in skills on first run.
5. ct doctor: verify all aqueduct.yaml-referenced skills exist in ~/.cistern/skills/.

Result: one canonical location at ~/.cistern/skills/, no copying, deterministic, always up to date. Skills managed in one place.

### From: manual

Implemented FINAL APPROACH: (1) session.go: add --add-dir ~/.cistern/skills to every claude invocation; (2) context.go: inject absolute skills.LocalPath(skill.Name) into CONTEXT.md <location> instead of sandbox-relative path; (3) runner.go: remove per-run copyFile loop and copyFile helper — no file copying needed; (4) doctor.go: check all skills (in-repo and external) against ~/.cistern/skills/, removing the skip for skill.Path \!= ''. Tests updated and all passing. Committed as 1317049.

### From: manual

Implemented: --add-dir ~/.cistern/skills in session.go, absolute path in context.go location, removed copyFile loop from runner.go, doctor.go now checks all skills. All tests pass.

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
    ct droplet pass ci-xilu9

**Recirculate (needs rework — send back upstream):**
    ct droplet recirculate ci-xilu9
    ct droplet recirculate ci-xilu9 --to implement

**Block (genuinely blocked, cannot proceed):**
    ct droplet block ci-xilu9

Add notes before signaling:
    ct droplet note ci-xilu9 "What you did / found"

The `ct` binary is on your PATH.
