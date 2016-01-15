package convert

import (
	"github.com/crockeo/go-tuner/synth"
	"io"
)

// Dealing with synth.RawDelayedNoteData from a MIDI file.
type MIDIArrangement struct{}

func (a MIDIArrangement) ReadNoteArrangement(reader io.Reader) ([]synth.RawDelayedNoteData, error) {
	return []synth.RawDelayedNoteData{}, nil
}

func (a MIDIArrangement) WriteNoteArrangement(writer io.Writer, notes []synth.RawDelayedNoteData) error {
	return nil
}
