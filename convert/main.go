package convert

import (
	"github.com/crockeo/go-tuner/filestore"
	"os"
	"path/filepath"
)

// Given a source path, destination path, and an arrangement type for both,
// convert the source file's format into the destination file's format.
func Convert(srcPath string, src filestore.ArrangementSource, dstPath string, dst filestore.ArrangementDestination) error {
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
	src, err := filestore.DecideFormat(filepath.Ext(srcPath))
	if err != nil {
		return err
	}

	dst, err := filestore.DecideFormat(filepath.Ext(dstPath))
	if err != nil {
		return err
	}

	return Convert(srcPath, src, dstPath, dst)
}
