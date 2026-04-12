# AGENTS.md

## Workflow

- Prefer `mise run <task>` for common project actions instead of invoking raw Go commands directly.
- Use `mise run format` for formatting, `mise run test` for tests, `mise run build` for builds, `mise run run` to launch the app, `mise run tidy` for module sync, `mise run generate-options` to refresh DOSBox enum values from the installed app, and `mise run check` for a quick format-plus-test pass.
- If a needed workflow is missing from `mise.toml`, add a task first when that will make future work more consistent.
- When changing bounded DOSBox settings in the UI, prefer regenerating option lists from the installed DOSBox Staging config instead of editing enum values by hand.

## Stack

- This project is a Go desktop app built with Fyne.
- DOSBox Staging integration is a core product concern, so changes should preserve the launch flow unless explicitly requested otherwise.
