module i3status-go

go 1.24.2

require (
	github.com/adrg/xdg v0.5.3
	github.com/davecgh/go-spew v1.1.1
	github.com/distatus/battery v0.11.0
	github.com/go-co-op/gocron/v2 v2.16.1
	github.com/hjson/hjson-go v3.3.0+incompatible
	github.com/mafik/pulseaudio v0.0.0-20240327130323-384e01075e6e
	github.com/shirou/gopsutil v3.21.11+incompatible
	go.i3wm.org/i3 v0.0.0-20190720062127-36e6ec85cc5a
)

require (
	github.com/BurntSushi/xgb v0.0.0-20210121224620-deaf085860bc // indirect
	github.com/BurntSushi/xgbutil v0.0.0-20190907113008-ad855c713046 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	howett.net/plist v1.0.1 // indirect
)

tool (
	i3status-go/cmd/battery-test
	i3status-go/cmd/i3status-go
	i3status-go/internal/lib
)
