package main

import (
	"flag"
	"log"
	"os"
)

var gozerDebug bool
var gozerVerbose bool

type GozerLog struct {
	log.Logger
}

var logger = GozerLog{
	*log.New(
		os.Stdout,
		"gozer: ",
		log.Lmicroseconds|log.Lshortfile)}

func (logger *GozerLog) Debugln(s ...string) {
	if gozerDebug {
		logger.Println(s)
	}
}

func (logger *GozerLog) Verboseln(s ...string) {
	if gozerVerbose || gozerDebug {
		logger.Println(s)
	}
}

func (logger *GozerLog) Error(err error, s ...string) error {
	logger.Println("ERROR: ", s, err)
	return err
}

func init() {
	flag.BoolVar(&gozerDebug, "debug", false, "Turn on debug trace")
	flag.BoolVar(&gozerVerbose, "verbose", false, "Turn on verbose log")
}
