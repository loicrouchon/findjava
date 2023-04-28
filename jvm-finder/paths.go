package main

import (
	"os"
	"os/user"
	"strings"
)

func resolvePath(path string) string {
	path = os.ExpandEnv(path)
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			die("Unable to resolve user home directory used in path %s: %s", path, err)
		}
		path = strings.Replace(path, "~", usr.HomeDir, 1)
	}
	return path
}
