package internal

import (
	"fmt"
	"net/netip"
	"syscall"

	"golang.org/x/sys/windows"
)

type MsgType int

const (
	NewIPAddrMsg MsgType = iota
	DelIPAddrMsg
	UpLinkMsg
	DownLinkMsg
)

var (
	iphlpapi             = windows.NewLazyDLL("Iphlpapi.dll")
	procNotifyAddrChange = iphlpapi.NewProc("NotifyAddrChange")
)

type MobilityMsg struct {
	MsgType       MsgType
	InterfaceName string
	Addr          netip.Addr
}

func RoutineMobilityMessageReceive(mmsgCh chan []MobilityMsg, ifacem *InterfaceManager) {
	fmt.Println("Listening for Mobility messages...")

	for {
		// wait ip address change
		r1, _, err := syscall.SyscallN(uintptr(procNotifyAddrChange.Addr()), 0, 0)

		if syscall.Errno(r1) != windows.NOERROR || err != windows.NOERROR {
			fmt.Println("Error receiving message: ", syscall.Errno(r1))
			continue
		}

		fmt.Println("IP Change Detected")

		// get now interfaces
		nowifaces := GetInterfaces()

		// compare interfaces
		// detect new interface
		for _, newiface := range *nowifaces {
			var isFound bool
			for _, iface := range ifacem.Interfaces {
				if iface.Name == newiface.Name {
					isFound = true
					break
				}
			}

			// detected new interface
			if !isFound {
				// create Mobility Message
				mmsg := MobilityMsg{MsgType: UpLinkMsg, InterfaceName: newiface.Name}
				mmsgCh <- []MobilityMsg{mmsg}
			}

		}

		// detect deleted interface
		for _, iface := range ifacem.Interfaces {
			var isFound bool
			for _, newiface := range *nowifaces {
				if iface.Name == newiface.Name {
					isFound = true
					break
				}
			}

			// detected deleted interface
			if !isFound {
				mmsg := MobilityMsg{MsgType: DownLinkMsg, InterfaceName: iface.Name}
				mmsgCh <- []MobilityMsg{mmsg}
			}
		}

		// compare IP Address
		for _, newiface := range *nowifaces {
			for _, iface := range ifacem.Interfaces {
				if iface.Name != newiface.Name {
					continue
				}

				// if iface is found. compare IP Address.
				// detect new IP Address
				// ipv4
				for _, newAddr := range newiface.IPv4Addr {
					var isFound bool
					for _, addr := range iface.IPv4Addr {
						if addr == newAddr {
							isFound = true
							break
						}
					}

					if !isFound {
						mmsg := MobilityMsg{MsgType: NewIPAddrMsg, InterfaceName: iface.Name, Addr: newAddr}
						mmsgCh <- []MobilityMsg{mmsg}
					}
				}

				// ipv6
				for _, newAddr := range newiface.IPv6Addr {
					var isFound bool
					for _, addr := range iface.IPv6Addr {
						if addr == newAddr {
							isFound = true
							break
						}
					}

					if !isFound {
						mmsg := MobilityMsg{MsgType: NewIPAddrMsg, InterfaceName: iface.Name, Addr: newAddr}
						mmsgCh <- []MobilityMsg{mmsg}

					}
				}

				// detect deleted IP Address
				// ipv4
				for _, addr := range iface.IPv4Addr {
					var isFound bool
					for _, newAddr := range newiface.IPv4Addr {
						if addr == newAddr {
							isFound = true
							break
						}
					}

					if !isFound {
						mmsg := MobilityMsg{MsgType: DelIPAddrMsg, InterfaceName: iface.Name, Addr: addr}
						mmsgCh <- []MobilityMsg{mmsg}
					}
				}

				// ipv6
				for _, addr := range iface.IPv6Addr {
					var isFound bool
					for _, newAddr := range newiface.IPv6Addr {
						if addr == newAddr {
							isFound = true
							break
						}
					}

					if !isFound {
						mmsg := MobilityMsg{MsgType: DelIPAddrMsg, InterfaceName: iface.Name, Addr: addr}
						mmsgCh <- []MobilityMsg{mmsg}
					}
				}
			}
		}

		// update interface manager
		ifacem.Interfaces = *nowifaces
	}
}
