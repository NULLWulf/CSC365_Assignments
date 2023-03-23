package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, os.Kill, syscall.SIGKILL)

	switch os.Args[1] {
	case "RUN1":
		log.Printf("Starting Assignment 1")
		RUNA1()
	case "RUN2_1":
		log.Printf("Starting Assignment 2: Loader")
		RUN2_1()
	case "RUN2_2":
		log.Printf("Starting Assignment 2: Application")
		RUN2_2()
	}

	// Wait for a termination signal
	<-signalChan
	log.Println("Termination signal received. Exiting...")
}
