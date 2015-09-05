package main

import (
	"fmt"
	"github.com/crockeo/go-tuner/filestore"
	"github.com/crockeo/go-tuner/server"
	"github.com/crockeo/go-tuner/synth"
	"os"
)

// Printing out help information for the user.
func printHelp() {
	fmt.Println("Usage:")
	fmt.Println(" go-tuner server")
	fmt.Println(" go-tuner file <file/path>")
	fmt.Println(" go-tuner visualize <file/path>")
}

// Handling er
func handleErrors(errChannel chan error) {
	for {
		err := <-errChannel
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

// The entry point to the application.
func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" {
		printHelp()
		return
	}

	if os.Args[1] == "server" {
		errChannel := make(chan error)
		noteChannel := make(chan synth.DelayedNoteData)

		go server.Start(errChannel, noteChannel)
		go handleErrors(errChannel)

		err := synth.StartDynamicSynth(noteChannel)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else if os.Args[1] == "file" {
		if len(os.Args) != 3 {
			printHelp()
			return
		}

		na, err := filestore.LoadNoteArrangement("res/songs/test.json")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = synth.StartStaticSynth(na)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else if os.Args[1] == "visualize" {
		fmt.Println("Visualize not yet implemented.")
	}
}
