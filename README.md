# fangcli

**A Command-Line Interface (CLI) to fangbreaker.zone for the EverQuest TLP Fangbreaker.**

> **Disclaimer:** This application is an independent project and is **not** associated with the fangbreaker.zone website, the developers of the site, or the video game EverQuest.

## Overview

`fangcli` provides a quick way to query and filter the zone information available on **fangbreaker.zone** directly from your terminal.

Whether you're looking for zones that offer a specific **bonus** (like experience or rare loot) or you need to find an appropriate zone within a certain **level range**, `fangcli` delivers the data without leaving your command line.

---

## Getting Started

### Installation

Go to https://github.com/andrewmriley/fangcli/releases and download the appropriate release for your operating system. `fanglci` does not write any files or read any files from your drive so it can be run from anywhere including removeable media.

Once installed, the `fangcli` command should be available in your terminal.

-----

## Usage and Arguments

Run the command with no arguments to see all available zone data.

```bash
fangcli
```

You can use the following flags to filter and sort the output:

| Argument       | Description                                                          | Values                                                                    |
|:---------------|:---------------------------------------------------------------------|:--------------------------------------------------------------------------|
| `-bonus`       | Limit the output to zones that list a specific bonus.                | `experience`, `AA`, `coin`, `loot`, `rare`, `skill`, `respawn`, `faction` |
| `-minlevel`    | Show zones **greater than or equal to** this level.                  | Integer                                                                   |
| `-maxlevel`    | Show zones **less than or equal to** this level.                     | Integer                                                                   |
| `-expansion`   | Limit the results to a specific expansion.                           | `classic`, `velious`, etc.                                                |
| `-zonetype`    | Limit the results to a specific type of zone.                        | `indoor`, `outdoor`                                                 |
| `-sortbylevel` | Change the sort direction of the zones by level. Default is **ASC**. | `ASC` or `DESC`                                                           |
| `-quiet`       | Hide extra output, focusing only on the zone list.                   | (Flag only)                                                               |

### Examples

#### Basic Running

Shows all zone data.

```bash
fangcli
```

#### Only Showing Experience Bonus Zones

```bash
fangcli -bonus experience
```

#### Show Zones at Least Level 10

```bash
fangcli -minlevel 10
```

#### Show Rare Bonus Zones Between Levels 10 and 30

```bash
fangcli -minlevel 10 -maxlevel 30 -bonus rare
```

#### Show all zones from the "Kunark" expansion, sorted high-to-low

```bash
fangcli -expansion kunark -sortbylevel DESC
```

-----

## Contributing

This project is a personal learning tool, but feedback and contributions are always welcome\!

If you find a bug or have a feature suggestion, please open an issue. If you'd like to contribute code, please check for existing issues or open a new one to discuss your proposed changes first.

-----

## License
This program is under the BSD 3-Clause License. See the LICENSE file for more information.