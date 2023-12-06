package main

import (
	"fmt"
	"os"
)

func createSocket(sock string) string {
	uid := os.Getuid()
	if uid == -1 {
		return ""
	}

	dir := fmt.Sprintf("/run/user/%d", uid)
	if _, err := os.Stat(dir); err != nil {
		return fmt.Sprintf("/run/%s", sock)
	} else {
		return fmt.Sprintf("/run/user/%d/%s", uid, sock)
	}
}
