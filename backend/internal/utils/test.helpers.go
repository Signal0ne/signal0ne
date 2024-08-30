package utils

import (
	"fmt"
	"net"
	"time"
)

func ConnectToSocket() net.Conn {
	var conn net.Conn = nil
	var err error
	for conn == nil {
		conn, err = net.DialTimeout("unix", "/var/run/d.sock", (60 * time.Second))
		if err != nil {
			fmt.Printf("Failed to establish connection, error: %s, retrying in 5 seconds\n",
				err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	return conn
}
