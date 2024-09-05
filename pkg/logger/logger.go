package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// var Logger zerolog.Logger

var A zerolog.Logger
var E zerolog.Logger
var I zerolog.Logger

func Init(filepath string) {

	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(filepath, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	var access_path = fmt.Sprintf("%v/access.log", filepath)
	var info_path = fmt.Sprintf("%v/info.log", filepath)
	var error_path = fmt.Sprintf("%v/error.log", filepath)

	access_file, err := os.OpenFile(access_path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	info_file, err := os.OpenFile(info_path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	error_file, err := os.OpenFile(error_path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	A = zerolog.New(zerolog.MultiLevelWriter(consoleWriter, access_file)).With().Timestamp().Logger()
	I = zerolog.New(zerolog.MultiLevelWriter(consoleWriter, info_file)).With().Timestamp().Logger()
	E = zerolog.New(zerolog.MultiLevelWriter(consoleWriter, error_file)).With().Timestamp().Logger()
}
