package filestore

import (
	"errors"
	"github.com/crockeo/go-tuner/synth"
	"io"
)

// Some source for a set of RawDelayedNoteData that can be written out using an
// ArrangementDestination
type ArrangementSource interface {
	ReadNoteArrangement(io.Reader) ([]synth.RawDelayedNoteData, error)
}

// A destination for some set of []synth.RawDelayedNoteData.
type ArrangementDestination interface {
	WriteNoteArrangement(io.Writer, []synth.RawDelayedNoteData) error
}

// Merging the functionality of ArrangementSource and ArrangementDestination
// into a single interface.
type ArrangementFormat interface {
	ReadNoteArrangement(io.Reader) ([]synth.RawDelayedNoteData, error)
	WriteNoteArrangement(io.Writer, []synth.RawDelayedNoteData) error
}

// Given a string representing a file extension, attempt to map it to an
// arrangement type.
func DecideFormat(extension string) (ArrangementFormat, error) {
	switch extension {
	case ".json":
		return JSONArrangement{}, nil
	case ".mid":
		return MIDIArrangement{}, nil
	case ".txt":
		return TextArrangement{}, nil
	default:
		return nil, errors.New("Undecidable extension.")
	}
}
