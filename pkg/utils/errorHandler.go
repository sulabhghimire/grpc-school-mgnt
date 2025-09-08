package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

func ErrorHandler(err error, msg string) error {
	// Capture the caller (1 level up from this function)
	_, file, line, ok := runtime.Caller(1)

	errorLogger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)

	if ok {
		errorLogger.Printf("%s:%d - %s | %v\n", file, line, msg, err)
	} else {
		errorLogger.Printf("%s | %v\n", msg, err)
	}

	return fmt.Errorf("error: %s", msg)
}
