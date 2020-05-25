package main

import (
	"log"
)

var (
	InfoLog  *log.Logger
	ErrorLog *log.Logger
)

func startLog() {
	InfoLog = log.New(InfoLogOutPut,
		"Info ",
		log.Ldate|log.Ltime|log.Lshortfile)

	ErrorLog = log.New(ErrorLogOutPut,
		"Error ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
