package internal

import (
	"fmt"
	"net"
	"net/netip"
)

// RegistInterfaces regist interfaces
func RegistInterfaces(ifaceManager *InterfaceManager, ifaces *[]Interface) {
	for _, iface := range *ifaces {
		ifaceManager.NewLink(iface.Name, iface.State)

		for _, addr := range iface.IPv4Addr {
			ifaceManager.NewIPAddr(iface.Name, addr)
		}
		for _, addr := range iface.IPv6Addr {
			ifaceManager.NewIPAddr(iface.Name, addr)
		}
	}
}

func GetInterfaces() *[]Interface {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Can't Retrieve interfaces")
	}

	var Interfaces []Interface

	for _, ifaceInfo := range interfaces {
		isUp := ifaceInfo.Flags&net.FlagUp != 0
		iface := Interface{Name: ifaceInfo.Name, State: isUp}

		addrs, _ := ifaceInfo.Addrs()

		for _, addr := range addrs {
			prefix, err := netip.ParsePrefix(addr.String())
			if err != nil {
				fmt.Println("Failed to ParsePrefix")
				continue
			}
			if prefix.Addr().Is4() {
				iface.IPv4Addr = append(iface.IPv4Addr, prefix.Addr())
			}
			if prefix.Addr().Is6() {
				iface.IPv6Addr = append(iface.IPv6Addr, prefix.Addr())
			}
		}

		Interfaces = append(Interfaces, iface)
	}

	return &Interfaces
}
