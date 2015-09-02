package main

import (
	"fmt"
	"github.com/crockeo/go-tuner/server"
	"os"
)

// Printing out help information for the user.
func printHelp() {
	fmt.Println("Usage:")
	fmt.Println(" go-tuner server")
	fmt.Println(" go-tuner file <file/path>")
	fmt.Println(" go-tuner visualize <file/path>")
}

// The entry point to the application.
func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" {
		printHelp()
		return
	}

	if os.Args[1] == "server" {
		err := server.Start(false)
		if err != nil {
			fmt.Println("Failed to start server: " + err.Error())
		}
	} else if os.Args[1] == "file" {
		fmt.Println("File not yet implemented.")
	} else if os.Args[1] == "visualize" {
		fmt.Println("Visualize not yet implemented.")
	}
}
