package filestore

import (
	"encoding/json"
	"github.com/crockeo/go-tuner/synth"
	"io"
)

// Dealing with synth.RawDelayedNoteData from a JSON file.
type JSONArrangement struct{}

func (a JSONArrangement) ReadNoteArrangement(reader io.Reader) ([]synth.RawDelayedNoteData, error) {
	dec := json.NewDecoder(reader)
	notes := []synth.RawDelayedNoteData{}

	err := dec.Decode(&notes)
	if err != nil {
		return []synth.RawDelayedNoteData{}, err
	}

	return notes, nil
}

func (a JSONArrangement) WriteNoteArrangement(writer io.Writer, notes []synth.RawDelayedNoteData) error {
	return json.NewEncoder(writer).Encode(notes)
}
