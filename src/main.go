package main

import (
	"fmt"

	"github.com/interfaec-alive-monitor-windows/internal"
)

func main() {
	var ifaceManager internal.InterfaceManager
	internal.RegistInterfaces(&ifaceManager, internal.GetInterfaces())
	ifaceManager.Print()

	mmsgCh := make(chan []internal.MobilityMsg, 64)
	defer close(mmsgCh)
	go internal.RoutineMobilityMessageReceive(mmsgCh, &ifaceManager)

	for {
		mmsgs := <-mmsgCh
		fmt.Println("-----")
		for _, mmsg := range mmsgs {
			switch mmsg.MsgType {
			case internal.NewIPAddrMsg:
				if mmsg.Addr.String() == "invalid IP" {
					continue
				}
				fmt.Print("NewIPAddr: ", mmsg.InterfaceName, " ", mmsg.Addr.String(), "\n")
			case internal.DelIPAddrMsg:
				if mmsg.Addr.String() == "invalid IP" {
					continue
				}
				fmt.Print("DelIPAddr: ", mmsg.InterfaceName, " ", mmsg.Addr.String(), "\n")
			case internal.UpLinkMsg:
				fmt.Print("UpLink: ", mmsg.InterfaceName, "\n")
			case internal.DownLinkMsg:
				fmt.Print("DownLink: ", mmsg.InterfaceName, "\n")
			}
			ifaceManager.Print()
		}
	}
}
