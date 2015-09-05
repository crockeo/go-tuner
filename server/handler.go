package server

import (
	"bytes"
	"encoding/json"
	"github.com/crockeo/go-tuner/synth"
)

//// Attempting to parse a given message into a DelayedNote. It will return
//// an error upon failure.
func ParseMessage(str string) (synth.DelayedNoteData, error) {
	var dnd synth.DelayedNoteData
	err := json.NewDecoder(bytes.NewBufferString(str)).Decode(&dnd)
	if err != nil {
		return synth.DelayedNoteData{}, nil
	}

	return dnd, nil
}

// Attempting to handle a message and send the parsed data over to the synth
// through a channel.
func HandleMessage(str string, noteChannel chan synth.DelayedNoteData) error {
	// If the command is a comment, just ignore it without sending anything into
	// the channel.
	if len(str) > 0 && str[0] == '/' {
		return nil
	}

	dnd, err := ParseMessage(str)
	if err != nil {
		return err
	}

	noteChannel <- dnd

	return nil
}
