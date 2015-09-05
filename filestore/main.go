package filestore

import (
	"encoding/json"
	"errors"
	"github.com/crockeo/go-tuner/synth"
	"os"
)

// Trying to load a synth.NoteArrangement from a file on the disk.
func LoadNoteArrangement(path string) (*synth.NoteArrangement, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, errors.New("File \"" + path + "\" does not exist.")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	na := new(synth.NoteArrangement)
	err = json.NewDecoder(file).Decode(na)
	if err != nil {
		return nil, err
	}

	return na, nil
}
