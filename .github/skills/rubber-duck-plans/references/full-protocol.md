# Rubber-Duck Plans — Convergence Loop Protocol

Every plan must pass a rubber-duck convergence loop before the user sees it.
No exceptions. No shortcuts. This is a quality gate, not a suggestion.

## When This Applies

This protocol applies whenever you are:

- Creating an implementation plan
- Designing an architecture or system
- Breaking a task into sub-tasks or todos
- Proposing a refactoring strategy
- Making any multi-step proposal that the user will review before execution
- Working in plan mode (`[[PLAN]]`)

**If it's a plan, it goes through the loop.**

## Forbidden Rationalizations

These are NOT valid reasons to skip the loop. If you catch yourself thinking any of these, **stop and run the loop anyway**:

- ❌ "This plan is small/trivial — rubber-duck would just slow things down"
- ❌ "The scope is well-understood, no design decisions needed"
- ❌ "Both items follow established patterns"
- ❌ "I can self-review for the same things rubber-duck would catch"
- ❌ "The user is in a hurry / low-latency response would be preferred"
- ❌ "Only N items in the plan, rubber-duck overhead exceeds benefit"
- ❌ "I rubber-ducked earlier in this session for a similar plan"

The protocol is **always-on**. There is no plan-size threshold below which it can be skipped. If a plan is truly trivial, the rubber-duck round will be quick — but you still have to do it.

## Pre-flight Gate Before exit_plan_mode

Before calling `exit_plan_mode`, you MUST be able to point to:

1. At least one prior `task` tool call in this session with `agent_type: "rubber-duck"` whose prompt included the current version of the plan
2. A terminal verdict of READY or READY WITH RISKS from that critique (or a later one if you iterated)
3. A rubber-duck-status section in `plan.md` summarizing the rounds and findings

If any of these are missing, do NOT call `exit_plan_mode`. Run the loop first.

## Recursion Guard

This protocol does **NOT** apply when you are:

- Currently acting as the rubber-duck reviewer (critiquing someone else's plan)
- Processing a convergence round (reviewing an updated plan within this loop)
- Performing any form of critique, review, audit, or validation of plans or
  outputs produced by this protocol itself

If the current task is itself a rubber-duck review, skip this protocol entirely.

## The Convergence Loop

### Step 1: Draft the Plan

Write the complete plan internally. Include:

- Problem statement and proposed approach
- Files to create/modify
- Key design decisions and trade-offs
- Todos with dependencies

Do NOT present this to the user yet.

### Step 2: Invoke the Rubber-Duck Agent

Use the `task` tool with `agent_type: "rubber-duck"` in sync mode. Provide a
**review packet** containing:

1. **Plan text** — the complete current plan
2. **Touched files** — files and components the plan will modify, with relevant
   code excerpts
3. **Assumptions** — key assumptions the plan makes about the codebase, APIs,
   or environment
4. **Changes since last round** — what was revised (omit for round 1)
5. **Carry-forward table** — ALL unresolved findings from prior rounds with
   current status (omit for round 1)
6. **Critique instruction** — include this verbatim:

> "Provide a comprehensive critique of this plan. Rate each finding as
> CRITICAL, HIGH, MEDIUM, or LOW. Carry forward all unresolved findings from
> prior rounds. State your terminal verdict: READY, READY WITH RISKS, or
> NOT READY."

### Step 3: Process Findings

For each finding in the rubber-duck response:

1. Assign an **ID** using the format `F{n}-R{round}` (e.g., F1-R1, F2-R1,
   F1-R2)
2. Record the **severity** (CRITICAL / HIGH / MEDIUM / LOW)
3. Record the **description**
4. Set initial **status** to OPEN

Then for each OPEN finding, do one of:

- **RESOLVE** it — revise the plan to address the finding
- **ACCEPT** it — provide a reasoned justification for why the finding
  doesn't apply or why the current approach is correct despite it. The rubber
  duck must agree in the next round for this to stick.
- **DEFER** it — only for MEDIUM or LOW severity. Document the risk and move
  on. CRITICAL and HIGH findings cannot be deferred.

### Step 4: Check Terminal State

Evaluate the rubber-duck's verdict:

| Terminal State | Criteria | Action |
|---------------|----------|--------|
| **READY** | Zero unresolved findings. Plan is sound. | → Go to Step 6 |
| **READY WITH RISKS** | No unresolved CRITICAL or HIGH findings. Remaining MEDIUM/LOW risks are documented and accepted. | → Go to Step 6 |
| **NOT READY** | CRITICAL or HIGH findings remain unresolved. | → Go to Step 5 |

### Step 5: Check for Stall

A **stall** occurs when no finding status has improved for one full round —
the same CRITICAL or HIGH findings remain unchanged across two consecutive
rounds.

**If stalled OR 3 rounds completed at NOT READY:**

Stop. Do NOT present the plan. Ask the user:

> "The rubber-duck review has not converged after [N] rounds. These concerns
> remain unresolved:
>
> [list unresolved CRITICAL/HIGH findings]
>
> How would you like to proceed?"

**If not stalled and under 3 rounds:** Return to Step 2 with the revised
plan and updated carry-forward table.

### Step 6: Present the Plan

Only after reaching READY or READY WITH RISKS:

1. **Out-of-scope-recording check.** Before presenting, verify: every item the
   plan names as out-of-scope, follow-up, deferred, or future work either
   (a) appears in the workspace's `issues.md` already, or (b) has a
   corresponding todo in the plan that records it there. If neither, the
   plan is NOT converged — return to Step 2 with the gap as a finding.
2. Present the final plan to the user
3. Append the **Convergence Summary** (see format below)

## Finding Schema

Every finding tracked through the loop must have:

```
ID:          F{n}-R{round}     (e.g., F1-R1, F3-R2)
Severity:    CRITICAL | HIGH | MEDIUM | LOW
Description: What the issue is
Status:      OPEN → RESOLVED | ACCEPTED | DEFERRED
```

**Status transitions:**

- `OPEN` → `RESOLVED`: Plan was revised to address the finding
- `OPEN` → `ACCEPTED`: Justification provided and rubber duck agreed
- `OPEN` → `DEFERRED`: Risk documented (MEDIUM/LOW only — CRITICAL/HIGH
  cannot be deferred)

## Carry-Forward Table

Starting from round 2, every rubber-duck invocation must include a
carry-forward table showing ALL findings from ALL prior rounds:

```
| ID     | Severity | Description              | Status   |
|--------|----------|--------------------------|----------|
| F1-R1  | HIGH     | Missing error handling    | RESOLVED |
| F2-R1  | MEDIUM   | Could add retry logic     | DEFERRED |
| F1-R2  | HIGH     | Race condition in cache   | OPEN     |
```

No finding may silently disappear. Every finding from every round must appear
in the carry-forward table until the loop exits.

## Convergence Summary Format

Append this to every plan presented to the user:

```markdown
## Convergence Summary

**Terminal State**: READY | **Rounds**: 2 | **Findings**: 5 total — 3 resolved, 1 accepted, 1 deferred

| ID | Severity | Description | Status |
|----|----------|-------------|--------|
| F1-R1 | HIGH | [description] | RESOLVED |
| F2-R1 | MEDIUM | [description] | ACCEPTED |
| F3-R1 | LOW | [description] | DEFERRED |
| F1-R2 | HIGH | [description] | RESOLVED |
| F2-R2 | MEDIUM | [description] | RESOLVED |
```

If terminal state is READY WITH RISKS, also include:

```markdown
### Accepted Risks
- F2-R1 (MEDIUM): [description] — [justification for accepting]
- F3-R1 (LOW): [description] — [justification for deferring]
```

## Rules

1. **Never skip the loop.** Even for "simple" or "obvious" plans. The whole
   point is catching blind spots you don't know you have.

2. **Never present a plan before convergence.** The user should never see a
   plan that hasn't reached READY or READY WITH RISKS. If you can't converge,
   ask the user — don't present a non-converged plan.

3. **Every finding must be addressed.** RESOLVED, ACCEPTED, or DEFERRED. No
   finding sits at OPEN when the loop exits.

4. **Track everything.** The carry-forward table is the source of truth. If
   a finding isn't in the table, it doesn't exist. If it IS in the table, it
   must have a status.

5. **No duplicate findings.** Each round should produce only NEW findings or
   status changes on existing ones. If the rubber duck raises the same issue
   again, it should reference the existing finding ID and explain why the
   resolution was insufficient.

6. **CRITICAL and HIGH cannot be deferred.** They must be RESOLVED or
   ACCEPTED (with rubber-duck agreement). Only MEDIUM and LOW can be deferred.

7. **Stall = ask user.** Never silently give up, never present a non-converged
   plan, never skip remaining CRITICAL/HIGH findings.

8. **Keep the summary compact.** The convergence summary is a quick-reference
   table, not a verbose narrative. The user should see terminal state, round
   count, and finding statuses at a glance.

9. **Out-of-scope items must be recorded.** Any plan that names items as
   out-of-scope, follow-up, deferred, or future work MUST record those items
   in the workspace's `issues.md` (or equivalent issue-tracking file) before
   convergence is declared. If the file cannot be edited safely (e.g. foreign
   uncommitted work per the concurrent-session-safety rule), the plan must
   include a deferred todo for the recording, gated to run as soon as it
   becomes safe.

   Items called "out of scope" that exist only in the plan and nowhere else
   are not actually tracked — they are silently discarded when the session
   ends. This rule closes that gap.

   The rubber-duck reviewer must call this out as a CRITICAL finding when
   it sees a plan with OOS items but no corresponding `issues.md` recording
   (either pre-existing or as a plan todo).

