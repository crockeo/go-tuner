package filestore

import (
	"errors"
	"github.com/crockeo/go-tuner/synth"
	"os"
	"path/filepath"
)

// Loading a synth.NoteArrangement form a file on disk.
func LoadNoteArrangement(path string) (*synth.NoteArrangement, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("Could not open \"" + path + "\".")
	}

	src, err := DecideFormat(filepath.Ext(path))
	if err != nil {
		return nil, err
	}

	rdnds, err := src.ReadNoteArrangement(file)
	if err != nil {
		return nil, err
	}

	na, err := synth.MakeNoteArrangement(rdnds)
	if err != nil {
		return nil, err
	}

	return na, nil
}
