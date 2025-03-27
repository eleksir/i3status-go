package lib

import (
	"fmt"
	"net"
	"os"
	"time"
)

// UpdateVPNStatus periodically update status of openvpn daemon.
func (c *MyConfig) UpdateVPNStatus() {
	var (
		vpnCStatus string
		tcpCheck   string
		vpnCheck   string
	)

	ticker := time.NewTicker(time.Second * 3)

	for range ticker.C {
		if c.VPNFileCheck() {
			if c.Vpn.UpColor == "" {
				vpnCheck = `⍋`
			} else {
				vpnCheck = fmt.Sprintf(`<span foreground="%s">⍋</span>`, c.Vpn.UpColor)
			}
		} else {
			if c.Vpn.DownColor == "" {
				vpnCheck = `⍒`
			} else {
				vpnCheck = fmt.Sprintf(`<span foreground="%s">⍒</span>`, c.Vpn.DownColor)
			}
		}

		if c.Vpn.TCPCheck.Enabled {
			if c.VPNTCPCheck() {
				if c.Vpn.UpColor == "" {
					tcpCheck = `✔`
				} else {
					tcpCheck = fmt.Sprintf(`<span foreground="%s">✔</span>`, c.Vpn.UpColor)
				}
			} else {
				if c.Vpn.DownColor == "" {
					tcpCheck = `✘`
				} else {
					tcpCheck = fmt.Sprintf(`<span foreground="%s">✘</span>`, c.Vpn.DownColor)
				}
			}

			vpnCStatus = fmt.Sprintf("VPN:%s:%s", vpnCheck, tcpCheck)
		} else {
			vpnCStatus = fmt.Sprintf("VPN:%s", vpnCheck)
		}

		if c.Values.VPNStatus != vpnCStatus {
			c.Values.VPNStatus = vpnCStatus
			c.Channels.UpdateReady <- true
		}
	}
}

// VPNTCPCheck intended to check arbitrary service inside vpn segment, to indicate that openvpn sevice not stoned.
func (c *MyConfig) VPNTCPCheck() bool {
	conn, err := net.DialTimeout(
		"tcp",
		fmt.Sprintf("%s:%d", c.Vpn.TCPCheck.Host, c.Vpn.TCPCheck.Port),
		time.Second*time.Duration(c.Vpn.TCPCheck.Timeout),
	)

	if err == nil {
		_ = conn.Close()
		return true
	}

	return false
}

// VPNFileCheck checks modification time of openvpn-status file.
func (c *MyConfig) VPNFileCheck() bool {
	var (
		fi  os.FileInfo
		err error
	)

	if fi, err = os.Stat(c.Vpn.StatusFile); err != nil {
		return false
	}

	if time.Since(fi.ModTime()).Seconds() > float64(c.Vpn.MtimeThreshold) {
		return false
	}

	return true
}
