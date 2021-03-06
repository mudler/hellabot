package hbot

import (
	"github.com/mudler/sendfd"
	"fmt"
	"net"
)

func (irc *IrcCon) StartUnixListener() {
	unaddr, err := net.ResolveUnixAddr("unix", irc.unixastr)
	if err != nil {
		panic(err)
	}
	list, err := net.ListenUnix("unix", unaddr)
	if err != nil {
		panic(err)
	}
	con, err := list.AcceptUnix()
	if err != nil {
		panic(err)
	}
	list.Close()

	fi, err := irc.con.(*net.TCPConn).File()
	if err != nil {
		panic(err)
	}

	err = sendfd.SendFD(con, fi)
	if err != nil {
		panic(err)
	}

	close(irc.Incoming)
	close(irc.outgoing)
}

// Attempt to hijack session previously running bot
func (irc *IrcCon) HijackSession() bool {
	unaddr, err := net.ResolveUnixAddr("unix", irc.unixastr)
	if err != nil {
		irc.Log(LWarning, "could not resolve unix socket")
		return false
	}

	con, err := net.DialUnix("unix", nil, unaddr)
	if err != nil {
		fmt.Println("Couldnt restablish connection, no prior bot.")
		fmt.Println(err)
		return false
	}

	ncon, err := sendfd.RecvFD(con)
	if err != nil {
		panic(err)
	}

	netcon, err := net.FileConn(ncon)
	if err != nil {
		panic(err)
	}

	irc.reconnect = true
	irc.con = netcon
	return true
}
