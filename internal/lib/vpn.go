package lib

import (
	"fmt"
	"net"
	"os"
	"time"
)

// UpdateVPNStatus periodically update status of openvpn daemon.
func (c MyConfig) UpdateVPNStatus() {
	var (
		vpnCStatus string
		tcpCheck   string
		vpnCheck   string
	)

	for {
		if c.VPNFileCheck() {
			if Conf.Vpn.UpColor == "" {
				vpnCheck = `⍋`
			} else {
				vpnCheck = fmt.Sprintf(`<span foreground="%s">⍋</span>`, c.Vpn.UpColor)
			}
		} else {
			if Conf.Vpn.DownColor == "" {
				vpnCheck = `⍒`
			} else {
				vpnCheck = fmt.Sprintf(`<span foreground="%s">⍒</span>`, c.Vpn.DownColor)
			}
		}

		if Conf.Vpn.TCPCheck.Enabled {
			if c.VPNTCPCheck() {
				if Conf.Vpn.UpColor == "" {
					tcpCheck = `✔`
				} else {
					tcpCheck = fmt.Sprintf(`<span foreground="%s">✔</span>`, c.Vpn.UpColor)
				}
			} else {
				if Conf.Vpn.DownColor == "" {
					tcpCheck = `✘`
				} else {
					tcpCheck = fmt.Sprintf(`<span foreground="%s">✘</span>`, c.Vpn.DownColor)
				}
			}

			vpnCStatus = fmt.Sprintf("VPN:%s:%s", vpnCheck, tcpCheck)
		} else {
			vpnCStatus = fmt.Sprintf("VPN:%s", vpnCheck)
		}

		if c.VPNStatus != vpnCStatus {
			c.VPNStatus = vpnCStatus
			c.UpdateReady <- true
		}

		time.Sleep(3 * time.Second)
	}
}

// VPNTCPCheck intended to check arbitrary service inside vpn segment, to indicate that openvpn sevice not stoned.
func (c MyConfig) VPNTCPCheck() bool {
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
func (c MyConfig) VPNFileCheck() bool {
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
