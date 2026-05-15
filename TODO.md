# TODO — Low-priority enhancements

Items deferred for future consideration.

## Output / UX

- **Structured output formats** (`--output json|table|yaml`): allow machine-readable output
  for commands such as `vm list`, `nic list`, `switch list`, `disk list`.
- **Color / icons in terminal output**: highlight VM states (running = green, off = red, etc.)
  using a lightweight ANSI helper.
- **Progress indicators** for long-running PowerShell operations (export, import, move).

## VM management

- **`vm clone` (linked clone / differencing disk)**: create a new VM backed by a differencing
  VHD derived from a parent disk without a full export/import cycle.
- **VM templates**: define a reusable VM template in YAML (generation, CPU, memory) and
  instantiate multiple VMs from it.
- **Checkpoint scheduling**: periodically create checkpoints via a cron-style flag or
  integration with Task Scheduler.

## Network

- **VLAN trunk mode**: `nic vlan --trunk` to configure `Set-VMNetworkAdapterVlan -Trunk`
  with allowed and native VLAN IDs.
- **NIC teaming / SR-IOV**: expose `Set-VMNetworkAdapter` SR-IOV and bandwidth management
  options.

## Storage

- **`disk optimize` / `disk compact`**: surface `Optimize-VHD` as a standalone subcommand.
- **`disk convert`**: wrap `Convert-VHD` for format changes (VHDX ↔ VHD, dynamic ↔ fixed).
- **`disk merge`**: wrap `Merge-VHD` for collapsing differencing disk chains.

## Host / infrastructure

- **Remote host support** (`--host <name>`): pass `-ComputerName` to applicable Hyper-V
  cmdlets so manahy can manage remote hosts.
- **Multi-host `list`**: extend `manahy list` to aggregate resource status across hosts
  defined in the YAML config.
- **`host replication`**: surface Hyper-V Replica configuration (`Enable-VMReplication`, etc.).

## Developer / CI

- **Integration / end-to-end tests**: tests that spin up a real Hyper-V environment
  (Windows runner in CI) to exercise full command paths.
- **Shell completion**: generate bash/zsh/fish/PowerShell completion scripts via
  `cobra.Command.GenBashCompletion`.
- **man page generation**: auto-generate man pages with `cobra/doc.GenManTree`.
