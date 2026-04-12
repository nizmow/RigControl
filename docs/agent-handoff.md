# Agent Handoff

## Purpose

RigControl is a Go/Fyne desktop app that builds DOSBox Staging machine configs from a GUI and launches DOSBox Staging with the generated config.

Current product shape:

- Pick a preset machine profile.
- Edit the profile in the UI.
- Preview the generated DOSBox config.
- Launch DOSBox Staging with a temporary config file.

## Current Architecture

Entry points:

- [cmd/rigcontrol/main.go](/Users/nizmow/Code/RigControl/cmd/rigcontrol/main.go): desktop app entry point.
- [cmd/generate-options/main.go](/Users/nizmow/Code/RigControl/cmd/generate-options/main.go): codegen tool for DOSBox enum-backed options.

Core packages:

- [internal/ui/app.go](/Users/nizmow/Code/RigControl/internal/ui/app.go): main Fyne window, form controls, preview, and launch flow.
- [internal/machine/profile.go](/Users/nizmow/Code/RigControl/internal/machine/profile.go): machine profile model used by the UI and config renderer.
- [internal/machine/presets.go](/Users/nizmow/Code/RigControl/internal/machine/presets.go): built-in starter presets.
- [internal/machine/options_generated.go](/Users/nizmow/Code/RigControl/internal/machine/options_generated.go): generated option lists for bounded DOSBox settings.
- [internal/dosbox/config.go](/Users/nizmow/Code/RigControl/internal/dosbox/config.go): DOSBox config rendering.
- [internal/dosbox/config_test.go](/Users/nizmow/Code/RigControl/internal/dosbox/config_test.go): basic renderer regression coverage.

## DOSBox Integration

The app currently hard-codes the DOSBox Staging executable path in [app.go](/Users/nizmow/Code/RigControl/internal/ui/app.go):

- `/Applications/DOSBox Staging.app/Contents/MacOS/dosbox`

Launch flow:

1. Build a `machine.Profile` from the form state.
2. Render a DOSBox config string.
3. Write a temporary `.conf` file.
4. Start DOSBox Staging with `-conf <tempfile>`.

Current limitation:

- The executable path is not yet configurable.

## Option Regeneration

Bounded DOSBox settings in the UI should come from the installed DOSBox Staging version, not from hand-maintained lists.

Current source of truth:

- `dosbox --printconf` returns the active config file path.
- That config file includes `Possible values:` comments for settings such as `machine`, `core`, `cputype`, `sbtype`, and `joysticktype`.

Generation flow:

1. Run `mise run generate-options`.
2. The generator reads the installed DOSBox config metadata.
3. It rewrites [options_generated.go](/Users/nizmow/Code/RigControl/internal/machine/options_generated.go).
4. The UI consumes the generated slices.

Do not hand-edit `options_generated.go`.

## Expected Workflows

Use `mise` tasks by default:

- `mise run run`
- `mise run build`
- `mise run format`
- `mise run test`
- `mise run check`
- `mise run tidy`
- `mise run generate-options`

## Current Product Decisions

- Launching DOSBox with the current config matters more than exporting files.
- The UI uses dropdowns only for settings that have documented bounded values.
- `cpu_cycles` remains free-form because DOSBox accepts numeric values and `max`.
- Presets are static starter templates, not persisted user data.

## Obvious Next Steps

- Make the DOSBox executable path configurable.
- Persist user-created profiles instead of only editing built-in presets in memory.
- Let users launch a game, folder, or disk image together with the selected machine config.
- Expand the generated option coverage beyond the current five bounded settings.
- Separate UI state management from widget construction once the form grows.

## Risks And Caveats

- DOSBox config keys and valid enum values can drift across versions, so regenerate options when the installed DOSBox version changes.
- The generator depends on a macOS app-bundle path today.
- Some current presets are only approximate “machine archetypes”, not historically strict hardware reproductions.
