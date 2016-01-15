package convert

import (
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
