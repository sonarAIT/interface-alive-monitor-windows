package main

import (
	"github.com/interfaec-alive-monitor-windows/internal"
)

func main() {
	var ifaceManager internal.InterfaceManager
	internal.RegistInterfaces(&ifaceManager)
	ifaceManager.Print()
}
