package convert

import (
	"errors"
	"os"
	"path/filepath"
)

// Given a string representing a file extension, attempt to map it to an
// arrangement type.
func DecideFormat(extension string) (ArrangementFormat, error) {
	switch extension {
	case ".json":
		return JSONArrangement{}, nil
	case ".midi":
		return MIDIArrangement{}, nil
	case ".txt":
		return TextArrangement{}, nil
	default:
		return nil, errors.New("Undecidable extension.")
	}
}

// Given a source path, destination path, and an arrangement type for both,
// convert the source file's format into the destination file's format.
func Convert(srcPath string, src ArrangementSource, dstPath string, dst ArrangementDestination) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	notes, err := src.ReadNoteArrangement(srcFile)
	if err != nil {
		return err
	}

	err = dst.WriteNoteArrangement(dstFile, notes)
	if err != nil {
		return err
	}

	return nil
}

// Similar to Convert, only that it tries to analyze the file extensions of the
// srcPath and dstPath to decide which arrangement types to use.
func ConvertAuto(srcPath string, dstPath string) error {
	src, err := DecideFormat(filepath.Ext(srcPath))
	if err != nil {
		return err
	}

	dst, err := DecideFormat(filepath.Ext(dstPath))
	if err != nil {
		return err
	}

	return Convert(srcPath, src, dstPath, dst)
}
