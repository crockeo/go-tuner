package main

import (
	"fmt"
	"github.com/crockeo/go-tuner/convert"
	"github.com/crockeo/go-tuner/filestore"
	"github.com/crockeo/go-tuner/server"
	"github.com/crockeo/go-tuner/synth"
	"github.com/crockeo/go-tuner/visualize"
	"os"
)

// Printing out help information for the user.
func printHelp() {
	fmt.Println("Usage:")
	fmt.Println(" go-tuner server")
	fmt.Println(" go-tuner file <file/path>")
	fmt.Println(" go-tuner visualize <file/path>")
	fmt.Println(" go-tuner convert <original/file/path> <new/file/path>")
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

	errChannel := make(chan error, 8)
	defer close(errChannel)
	go handleErrors(errChannel)

	if os.Args[1] == "server" {
		noteChannel := make(chan synth.DelayedNoteData, 32)
		defer close(noteChannel)
		go server.Start(errChannel, noteChannel)

		err := synth.StartSynth(noteChannel)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else if os.Args[1] == "file" {
		if len(os.Args) != 3 {
			printHelp()
			return
		}

		na, err := filestore.LoadNoteArrangement(os.Args[2])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if err = synth.StartSynthWith(na, make(chan synth.DelayedNoteData), true); err != nil {
			fmt.Println(err.Error())
		}
	} else if os.Args[1] == "visualize" {
		if len(os.Args) != 3 {
			printHelp()
			return
		}

		quitChannel := make(chan bool)
		defer close(quitChannel)

		noteChannel := make(chan synth.DelayedNoteData, 32)
		go synth.StartSynthAsync(noteChannel, quitChannel, errChannel)

		na, err := filestore.LoadNoteArrangement(os.Args[2])
		if err != nil {
			fmt.Println("Could not load song: " + err.Error())
		}

		if err = visualize.RunVisualization(na, noteChannel); err != nil {
			fmt.Println("Visualize error: " + err.Error())
		}

		quitChannel <- true
	} else if os.Args[1] == "convert" {
		if len(os.Args) != 4 {
			printHelp()
			return
		}

		if err := convert.ConvertAuto(os.Args[2], os.Args[3]); err != nil {
			fmt.Println("Failed to convert files: " + err.Error())
		}
	} else {
		printHelp()
	}
}
