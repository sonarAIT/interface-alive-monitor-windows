package internal

import (
	"fmt"
	"net/netip"
)

type MsgType int

const (
	NewIPAddrMsg MsgType = iota
	DelIPAddrMsg
	UpLinkMsg
	DownLinkMsg
)

type MobilityMsg struct {
	MsgType       MsgType
	InterfaceName string
	Addr          netip.Addr
}

func createMovementDetecterSocket() (int, error) {
}

func handleMobilityMessage(buf []byte) []NetlinkMsg {
}

func parseAddrMessage(buf []byte) (string, netip.Addr) {
}

func parseLinkMessage(buf []byte) (string, bool) {
}

func RoutineMobilityMessageReceive(nlmsgCh chan []MobilityMsg) {
	// fd, err := createMovementDetecterSocket()
	// if err != nil {
	// 	return
	// }
	// defer syscall.Close(fd)

	fmt.Println("Listening for Mobility messages...")

	// loop of receive mobility message
	// buf := make([]byte, 4096)
	for {
		// receive message
		// n, _, err := syscall.Recvfrom(fd, buf, 0)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Error receiving message: %v\n", err)
		// 	continue
		// }

		// handle mobility message
		// nlmsg := handleMobilityMessage(buf[:n])
		// nlmsgCh <- nlmsg
	}
}
