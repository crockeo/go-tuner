package convert

import (
	"github.com/crockeo/go-tuner/synth"
	"io"
)

// Dealing with synth.RawDelayedNoteData from a JSON file.
type JSONArrangement struct{}

func (a JSONArrangement) ReadNoteArrangement(reader io.Reader) ([]synth.RawDelayedNoteData, error) {
	return []synth.RawDelayedNoteData{}, nil
}

func (a JSONArrangement) WriteNoteArrangement(writer io.Writer, notes []synth.RawDelayedNoteData) error {
	return nil
}
