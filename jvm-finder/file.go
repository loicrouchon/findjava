package main

import "os"

func closeFile(file *os.File) {
	func(file *os.File) {
		err := file.Close()
		if err != nil {
			dierr(err)
		}
	}(file)
}
