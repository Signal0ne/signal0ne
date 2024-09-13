package utils

import (
	"fmt"
	"net"
	"signal0ne/cmd/config"
	"time"
)

var NUMBER_OF_RETRIES = 6
var RETRY_INTERVAL = 5

func ConnectToSocket() net.Conn {
	var conn net.Conn = nil
	var err error
	var cfg = config.GetInstance()
	if cfg.Server.Mode != "prod" {
		NUMBER_OF_RETRIES = 3
		RETRY_INTERVAL = 1
	}

	retries := 0
	for conn == nil && retries < NUMBER_OF_RETRIES {
		conn, err = net.DialTimeout("unix", cfg.IPCSocket, (60 * time.Second))
		if err != nil || conn == nil {
			fmt.Printf("Failed to establish connection, error: %s, retrying in %d seconds\n",
				err, RETRY_INTERVAL)
			time.Sleep(time.Duration(RETRY_INTERVAL) * time.Second)
		} else {
			break
		}
		retries++
	}
	return conn
}
