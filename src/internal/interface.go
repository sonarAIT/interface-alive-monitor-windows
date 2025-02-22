package internal

import (
	"fmt"
	"net/netip"
	"sync"
)

// Interface is represent machine network interface
// Lists head IP Addr is Primary IP Address.
type Interface struct {
	Name     string
	IPv4Addr []netip.Addr
	IPv6Addr []netip.Addr
	State    bool
}

// InterfaceManager is interface list holder of machine
type InterfaceManager struct {
	Interfaces []Interface
	sync.RWMutex
}

// NewLink add network interface to interface manager
func (ifacem *InterfaceManager) NewLink(ifaceName string, state bool) {
	// lock
	ifacem.Lock()
	defer ifacem.Unlock()

	// add iface
	newIface := Interface{Name: ifaceName, State: state}
	ifacem.Interfaces = append(ifacem.Interfaces, newIface)
}

// NewLink delete network interface from interface manager
func (ifacem *InterfaceManager) DelLink(ifaceName string) {
	// lock
	ifacem.Lock()
	defer ifacem.Unlock()

	// remove iface
	removeIface(&ifacem.Interfaces, ifaceName)
}

// searchInterfaceIdxFromIfaceName search interface index from interface name
func (ifacem *InterfaceManager) searchInterfaceIdxFromIfaceName(ifaceName string) int {
	for i, iface := range ifacem.Interfaces {
		if iface.Name == ifaceName {
			return i
		}
	}
	return -1
}

// NewIPAddr add IP Address to IP Address list in Interface
func (ifacem *InterfaceManager) NewIPAddr(ifaceName string, addr netip.Addr) {
	// lock
	ifacem.Lock()
	defer ifacem.Unlock()

	// search iface idx
	ifaceIdx := ifacem.searchInterfaceIdxFromIfaceName(ifaceName)

	// if iface is found. add addr.
	if ifaceIdx != -1 {
		if addr.Is4() {
			ifacem.Interfaces[ifaceIdx].IPv4Addr = append(ifacem.Interfaces[ifaceIdx].IPv4Addr, addr)
		} else if addr.Is6() {
			ifacem.Interfaces[ifaceIdx].IPv6Addr = append(ifacem.Interfaces[ifaceIdx].IPv6Addr, addr)
		}
		return
	}

	// if iface is not found. add iface with addr.
	newIface := Interface{Name: ifaceName}
	if addr.Is4() {
		newIface.IPv4Addr = []netip.Addr{addr}
	} else {
		newIface.IPv6Addr = []netip.Addr{addr}
	}
	ifacem.Interfaces = append(ifacem.Interfaces, newIface)
}

// NewIPAddr delete IP Address from IP Address list in Interface
func (ifacem *InterfaceManager) DelIPAddr(ifaceName string, addr netip.Addr) {
	// lock
	ifacem.Lock()
	defer ifacem.Unlock()

	// search iface
	ifaceIdx := ifacem.searchInterfaceIdxFromIfaceName(ifaceName)

	// if not found. return.
	if ifaceIdx == -1 {
		return
	}

	// remove addr
	if addr.Is4() {
		removeAddr(&ifacem.Interfaces[ifaceIdx].IPv4Addr, addr)
	} else {
		removeAddr(&ifacem.Interfaces[ifaceIdx].IPv6Addr, addr)
	}
}

// UpLink sets the state of the interface in the interface manager to up
func (ifacem *InterfaceManager) UpLink(ifaceName string) {
	// lock
	ifacem.Lock()
	defer ifacem.Unlock()

	// search iface
	ifaceIdx := ifacem.searchInterfaceIdxFromIfaceName(ifaceName)

	// if not found. return.
	if ifaceIdx == -1 {
		return
	}

	// set state
	ifacem.Interfaces[ifaceIdx].State = true
}

// DownLink sets the state of the interface in the interface manager to down
func (ifacem *InterfaceManager) DownLink(ifaceName string) {
	// lock
	ifacem.Lock()
	defer ifacem.Unlock()

	// search iface
	ifaceIdx := ifacem.searchInterfaceIdxFromIfaceName(ifaceName)

	// if not found. return.
	if ifaceIdx == -1 {
		return
	}

	// set state
	ifacem.Interfaces[ifaceIdx].State = false
}

// Print print interface list
func (ifacem *InterfaceManager) Print() {
	fmt.Println("-----")
	for _, iface := range ifacem.Interfaces {
		fmt.Println("{")
		fmt.Println("\tInterfaceName: ", iface.Name)

		if len(iface.IPv4Addr) == 0 {
			fmt.Println("\t Primary IPv4: ")
			fmt.Println("\t Secondary IPv4: []")
		} else {
			fmt.Println("\tPrimary IPv4: ", iface.IPv4Addr[0].String())
			fmt.Print("\tSecondary IPv4: [")
			for _, addr := range iface.IPv4Addr[1:] {
				fmt.Print(addr.String(), ",")
			}
			fmt.Println("]")
		}

		if len(iface.IPv6Addr) == 0 {
			fmt.Println("\t Primary IPv6: ")
			fmt.Println("\t Secondary IPv6: []")
		} else {
			fmt.Println("\tPrimary IPv6: ", iface.IPv6Addr[0].String())
			fmt.Print("\tSecondary IPv6: [")
			for _, addr := range iface.IPv6Addr[1:] {
				fmt.Print(addr.String(), ",")
			}
			fmt.Println("]")
		}

		fmt.Println("\tState: ", iface.State)
		fmt.Println("},")
	}
}

// removeAddr remove IP Address from netip.Addr slice
func removeAddr(slice *[]netip.Addr, addr netip.Addr) {
	for i, v := range *slice {
		if v == addr {
			*slice = (*slice)[:i+copy((*slice)[i:], (*slice)[i+1:])]
			return
		}
	}
}

// removeIface remove Interface from Interface slice
func removeIface(slice *[]Interface, ifaceName string) {
	for i, v := range *slice {
		if v.Name == ifaceName {
			*slice = (*slice)[:i+copy((*slice)[i:], (*slice)[i+1:])]
		}
	}
}
