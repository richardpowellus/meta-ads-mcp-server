---
name: cpd
description: "Clean, commit, push, and deploy this repo. Use for cpd, deploy, commit push deploy, or push and deploy."
---

# CPD - Clean, Commit, Push, Deploy

Use only in workspaces that have a real deploy path.

1. Do the CCP workflow: clean safe temp artifacts, inspect for secrets, stage only intended changes, validate, commit with the Copilot co-author trailer, and push without force.
2. Identify the repo's existing deploy command from project files/docs. Prefer, in order: explicit user instruction; `package.json` deploy script; `deploy.ps1`/`Deploy.ps1`; `deploy.sh`; documented deploy command in repo docs. Do not invent a deploy path.
3. If the user requested a target or bump (`cpd patch`, `cpd minor`, `cpd ios`, etc.), use only commands the repo already supports. If unsupported, stop and say so.
4. Run the deploy only after the push succeeds. Preserve deploy output and read failures literally.
5. Never deploy with unrelated dirty files, unpushed commits, missing required secrets, or an unknown deploy command. Never force-push.
6. Report commit hash, push target, deploy command, and deploy result.
