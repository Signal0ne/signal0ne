package utils

import (
	"fmt"
	"net"
	"signal0ne/cmd/config"
	"time"
)

func ConnectToSocket() net.Conn {
	var conn net.Conn = nil
	var err error
	var cfg = config.GetInstance()
	for conn == nil {
		conn, err = net.DialTimeout("unix", cfg.IPCSocket, (60 * time.Second))
		if err != nil || conn == nil {
			fmt.Printf("Failed to establish connection, error: %s, retrying in 5 seconds\n",
				err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	return conn
}
