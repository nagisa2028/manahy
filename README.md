# manahy

manahy is a management tool for Hyper-V, operated from the command line.

- **Windows only** — requires Hyper-V to be enabled and PowerShell available
- Membership in the **Hyper-V Administrators** group is required (or run as Administrator)

## Setup

### Build from source

```
> go mod download
> go build
```

### Character encoding (non-ASCII environments)

If PowerShell output appears garbled (e.g. Japanese locale), run the following command once per console session to switch the code page to UTF-8:

```
> chcp 65001
```

To make this permanent, add it to your PowerShell profile (`$PROFILE`) or change the system locale to Unicode via Windows regional settings.

## Environment check

Before first use, verify that Hyper-V is enabled and all required cmdlets are available:

```
> manahy host check
```

Subcommands for per-category detail:

```
> manahy host check vm
> manahy host check disk
> manahy host check network
```

## Global flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--config` | `-c` | `manahy.yaml` | Config file path (used by `build` / `remove` and for disk alias resolution) |
| `--help` | `-h` | | Show help |

The `-c` flag is inherited by all subcommands:

```
> manahy -c infra.yaml build
> manahy -c infra.yaml disk info test-disk
```

## Commands

### vm — virtual machine management

```
> manahy vm list
> manahy vm state <name>
> manahy vm create -n <name> -m <memory> -p <path> [options]
> manahy vm remove <name> [name...]
> manahy vm rename <name> --new-name <new>
> manahy vm start   <name> [name...]
> manahy vm shutdown <name> [name...]
> manahy vm destroy  <name> [name...]
> manahy vm save     <name> [name...]
> manahy vm suspend  <name> [name...]
> manahy vm resume   <name> [name...]
> manahy vm restart  <name> [name...]
> manahy vm connect  <name>
> manahy vm export   <name> --path <dir>
> manahy vm import   <path>
> manahy vm move     <name> --path <dir>
> manahy vm copy     <name> --new-name <new> --path <dir>
> manahy vm info     <name>
> manahy vm measure  <name>
> manahy vm integration list   <name>
> manahy vm integration enable  <name> --name <service>
> manahy vm integration disable <name> --name <service>
```

`vm create` options:

| Flag | Short | Description |
|---|---|---|
| `--name` | `-n` | VM name (required) |
| `--memory` | `-m` | Memory size, e.g. `512MB`, `2GB` (required) |
| `--path` | `-p` | Directory to store VM files (required) |
| `--generation` | `-g` | VM generation: `1` or `2` (default: `1`) |
| `--vcpus` | `-v` | vCPU count (default: `1`) |
| `--nested` | | Enable nested virtualization |
| `--nodynamic` | | Disable dynamic memory |
| `--disk` | `-d` | VHD path or config alias |
| `--network` | `-s` | Virtual switch name |
| `--image` | `-i` | ISO image path |

### disk — virtual disk (VHD) management

```
> manahy disk create --path <path> --size <size> --type <type>
> manahy disk remove <path>
> manahy disk info   <path>
> manahy disk resize <path> --size <size>
> manahy disk optimize <path>
> manahy disk convert  <path> --dest <dest>
> manahy disk mount    <path>
> manahy disk dismount <path>
> manahy disk merge    <path>
```

Disk types: `dynamic`, `fixed`, `differencing`

### switch — virtual switch management

```
> manahy switch list
> manahy switch create -n <name> -t <type> [--external-interface <adapter>]
> manahy switch remove <name>
> manahy switch rename <name> --new-name <new>
> manahy switch configure type    <name> --type <type>
> manahy switch configure adapter <name> --adapter <adapter>
```

Switch types: `external`, `internal`, `private`

### checkpoint — VM checkpoint management

```
> manahy checkpoint create <vm> [--name <name>]
> manahy checkpoint list   <vm>
> manahy checkpoint restore <vm> --name <name>
> manahy checkpoint remove  <vm> --name <name>
> manahy checkpoint rename  <vm> --name <name> --new-name <new>
> manahy checkpoint export  <vm> --name <name> --path <dir>
```

### dvd — VM DVD drive management

```
> manahy dvd list   <vm>
> manahy dvd add    <vm>
> manahy dvd remove <vm>
> manahy dvd set    <vm> --image <iso>
> manahy dvd eject  <vm>
```

### nic — VM network adapter management

```
> manahy nic list    <vm>
> manahy nic add     <vm> [--name <name>] [--switch <switch>]
> manahy nic remove  <vm> --name <name>
> manahy nic connect <vm> --name <name> --switch <switch>
```

### host — host information and management

```
> manahy host show
> manahy host check [vm|disk|network]
> manahy host member list
> manahy host member add    <user> [user...]
> manahy host member remove <user> [user...]
> manahy host storage list
```

## Config file (manahy.yaml)

`build` and `remove` operate on all resources defined in the config file at once.
The map key is the resource name; no separate `name:` field is needed.
Disk map keys can be used as **aliases** in VM `disks:` lists.

```yaml
disks:
  boot-disk:            # alias used in vm disks list below
    path: C:\VMs\boot.vhd
    size: 50GB
    type: dynamic
  data-disk:
    path: C:\VMs\data.vhd
    size: 100GB
    type: dynamic
  existing-disk:
    path: C:\VMs\existing.vhd
    import: true        # skip creation; treat as already existing

networks:
  internal-net:
    type: internal
  external-net:
    type: external
    external-interface: Ethernet
    allow-management-os: true

vms:
  web:
    count: 3            # creates web1, web2, web3
    generation: 2
    memory:
      size: 2GB
      dynamic: true
    cpu:
      thread: 2
      nested: false
    path: C:\VMs
    image: C:\ISOs\debian.iso
    disks:
      - boot-disk       # resolved to C:\VMs\boot.vhd
      - data-disk
    networks:
      - internal-net
```

#### Disk behaviour with `count > 1`

When `count` is greater than 1, each VM instance gets its **own numbered disk copy**.
The numeric suffix is inserted before the file extension:

| alias | `count` | disks created |
|-------|---------|---------------|
| `boot-disk` (`C:\VMs\boot.vhd`) | `3` | `boot1.vhd`, `boot2.vhd`, `boot3.vhd` |

`remove` deletes the same numbered copies; the base path (e.g. `boot.vhd`) is never touched.

Disks marked `import: true` are treated as **read-only references** and are shared across all instances without copying or deletion.

### build / remove

```
> manahy build              # uses manahy.yaml
> manahy -c infra.yaml build
> manahy build --dry-run    # preview what would be created
> manahy remove --dry-run   # preview what would be removed
> manahy remove             # remove all resources defined in the config
```

`--dry-run` output example (with `count: 3`):

```
[dry-run] build
  disk      boot-disk1                create    C:\VMs\boot1.vhd (50GB, dynamic)
  vm        web1                      create    gen2, 2GB, 2 vCPU
  disk      boot-disk2                create    C:\VMs\boot2.vhd (50GB, dynamic)
  vm        web2                      create    gen2, 2GB, 2 vCPU
  disk      boot-disk3                create    C:\VMs\boot3.vhd (50GB, dynamic)
  vm        web3                      create    gen2, 2GB, 2 vCPU
  network   internal-net              create    internal
```
