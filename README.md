# i3status-go

It is i3bar text generator.

## What it can

* Application launcher buttons
* Memory statistics
* LA, last 5 minutes
* Show battery charge
* Show average CPU cores temperature
* Network interfaces status
* OpenVPN status (including tcp checks)
* PulseAudio volume indicator, can adjust master volume too
* Clock
* Cron jobs (intended to use for show periodic desktop notifications but not limited to)
* Show output of one-shot system command

## How to build it

This application was tested to compile with go v1.24.

All you need is go lang compiler and gnu make utility. Invoke

```bash
make
```

that will produce binary *i3status-go*. You should place it to some dir intended for storing user binaries (~/.local/bin,
 ~/bin, or something else) and edit i3wm config according to
[manual](https://i3wm.org/docs/userguide.html#_configuring_i3bar) after editing i3wm config either run **i3status-go** by
hands and kill it or copy **i3status-go-example.json** to **$XDG_CONFIG_HOME/i3status-go.json** and adjust it to your
habbit.
