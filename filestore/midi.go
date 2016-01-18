package filestore

import (
	"errors"
	"github.com/crockeo/go-tuner/filestore/midi"
	"github.com/crockeo/go-tuner/synth"
	"io"
)

// Dealing with synth.RawDelayedNoteData from a MIDI file.
type MIDIArrangement struct{}

func (a MIDIArrangement) ReadNoteArrangement(reader io.Reader) ([]synth.RawDelayedNoteData, error) {
	m, err := midi.Read(reader)
	if err != nil {
		return []synth.RawDelayedNoteData{}, err
	}

	switch m.Header.Format {
	case 0:
	case 1:
	case 2:
	default:
		return []synth.RawDelayedNoteData{}, errors.New("Unexpected format.")
	}

	return []synth.RawDelayedNoteData{}, nil
}

func (a MIDIArrangement) WriteNoteArrangement(writer io.Writer, notes []synth.RawDelayedNoteData) error {
	return errors.New("WriteNoteArrangement not yet implemented.")
}
