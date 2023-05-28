package utils

import (
	"jvm-finder/internal/log"
	"os"
)

func WriteFile(name string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}

func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Die(err)
	}
}
