---
name: rubber-duck-plans
description: "Plan quality gate: before showing implementation/refactor plans, run rubber-duck review until READY or READY WITH RISKS."
---

# Rubber-Duck Plans

Use for any implementation plan, architecture decision, task breakdown, or refactor proposal before presenting it to the user.

Core loop:
1. Draft the plan internally, including approach, touched files, assumptions, dependencies, and risks.
2. Invoke `task` with `agent_type: "rubber-duck"` in sync mode. Include the full plan, relevant file context, assumptions, changes since prior round, and unresolved findings.
3. Ask the reviewer to rate findings CRITICAL/HIGH/MEDIUM/LOW and return `READY`, `READY WITH RISKS`, or `NOT READY`.
4. Track every finding by ID, severity, description, and status. Resolve or accept all CRITICAL/HIGH findings; do not silently drop carry-forward findings.
5. Iterate until `READY` or `READY WITH RISKS`. If stalled after 3 rounds or unresolved CRITICAL/HIGH findings remain, ask the user how to proceed instead of presenting the plan.
6. Present the final plan with a short convergence summary and accepted risks.

Read `references/full-protocol.md` when the plan is multi-phase, high-risk, cross-workspace, financial, production-affecting, or when a rubber-duck round returns anything other than READY.
