package server

import (
	"github.com/crockeo/go-tuner/synth"
)

// Attempting to handle a message and send the parsed data over to the synth
// through a channel.
func HandleMessage(str string, noteChannel chan synth.DelayedNoteData) error {
	// If the command is a comment, just ignore it without sending anything into
	// the channel.
	if len(str) > 0 && str[0] == '/' {
		return nil
	}

	dnd, err := synth.ParseMessage(str)
	if err != nil {
		return err
	}

	noteChannel <- dnd

	return nil
}
