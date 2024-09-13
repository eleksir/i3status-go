package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

// VPNStatus status string for i3bar, related to vpn.
var VPNStatus string

// UpdateVPNStatus periodically update status of openvpn daemon.
func UpdateVPNStatus() {
	var (
		vpnCStatus string
		tcpCheck   string
		vpnCheck   string
	)

	for {
		if VPNFileCheck() {
			if Conf.Vpn.UpColor == "" {
				vpnCheck = `⍋`
			} else {
				vpnCheck = fmt.Sprintf(`<span foreground="%s">⍋</span>`, Conf.Vpn.UpColor)
			}
		} else {
			if Conf.Vpn.DownColor == "" {
				vpnCheck = `⍒`
			} else {
				vpnCheck = fmt.Sprintf(`<span foreground="%s">⍒</span>`, Conf.Vpn.DownColor)
			}
		}

		if Conf.Vpn.TCPCheck.Enabled {
			if VPNTCPCheck() {
				if Conf.Vpn.UpColor == "" {
					tcpCheck = `✔`
				} else {
					tcpCheck = fmt.Sprintf(`<span foreground="%s">✔</span>`, Conf.Vpn.UpColor)
				}
			} else {
				if Conf.Vpn.DownColor == "" {
					tcpCheck = `✘`
				} else {
					tcpCheck = fmt.Sprintf(`<span foreground="%s">✘</span>`, Conf.Vpn.DownColor)
				}
			}

			vpnCStatus = fmt.Sprintf("VPN:%s:%s", vpnCheck, tcpCheck)
		} else {
			vpnCStatus = fmt.Sprintf("VPN:%s", vpnCheck)
		}

		if VPNStatus != vpnCStatus {
			VPNStatus = vpnCStatus
			UpdateReady <- true
		}

		time.Sleep(3 * time.Second)
	}
}

// VPNTCPCheck intended to check arbitrary service inside vpn segment, to indicate that openvpn sevice not stoned.
func VPNTCPCheck() bool {
	conn, err := net.DialTimeout(
		"tcp",
		fmt.Sprintf("%s:%d", Conf.Vpn.TCPCheck.Host, Conf.Vpn.TCPCheck.Port),
		time.Second*time.Duration(Conf.Vpn.TCPCheck.Timeout),
	)

	if err == nil {
		_ = conn.Close()
		return true
	}

	return false
}

// VPNFileCheck checks modification time of openvpn-status file.
func VPNFileCheck() bool {
	var (
		fi  os.FileInfo
		err error
	)

	if fi, err = os.Stat(Conf.Vpn.StatusFile); err != nil {
		return false
	}

	if time.Since(fi.ModTime()).Seconds() > float64(Conf.Vpn.MtimeThreshold) {
		return false
	}

	return true
}
