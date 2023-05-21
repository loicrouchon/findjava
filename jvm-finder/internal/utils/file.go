package utils

import (
	"jvm-finder/internal/log"
	"os"
)

func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Die(err)
	}
}
