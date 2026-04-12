# RigControl

RigControl is a machine manager for DOSBox Staging. The app uses a desktop GUI to assemble machine profiles such as a 286 with EGA or a fast Pentium, preview the resulting DOSBox config, and launch DOSBox Staging with that generated config.

## Stack

- Go
- Fyne
- `mise` for toolchain and task management

## Getting started

```sh
mise install
mise run tidy
mise run run
```

## Project Notes

- Future handoff and architecture notes live in [docs/agent-handoff.md](/Users/nizmow/Code/RigControl/docs/agent-handoff.md).
