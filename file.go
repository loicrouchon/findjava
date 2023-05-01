package main

import "os"

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		die(err)
	}
}
