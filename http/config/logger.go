package config

import (
	"fmt"
	"io"
	"log"
	"os"
)

func createdir(dir string) (bool, error) {
	_, err := os.Stat(dir)

	if err == nil {
		//directory exists
		return true, nil
	}

	err2 := os.MkdirAll(dir, 0755)
	if err2 != nil {
		return false, err2
	}

	return true, nil
}

// SetLogger 设置logger
func SetLogger(logger *log.Logger) {
	r, err := createdir(LogPath)
	if r == false {
		panic(err)
	}
	file, _ := os.OpenFile(LogPath+"/hanabi.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	writers := []io.Writer{
		file,
		os.Stdout,
	}
	logger.SetOutput(io.MultiWriter(writers...))
	logger.SetPrefix(fmt.Sprintf("[%s]", NAME))
	logger.SetFlags(log.Ldate | log.Ltime)
}
