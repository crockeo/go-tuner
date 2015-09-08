package convert

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crockeo/go-tuner/synth"
	"io"
	"os"
	"strings"
)

// Trying to parse DelayedNoteData out of a line.
func parseNoteData(line string) (synth.RawDelayedNoteData, error) {
	var delay float32
	var note string
	var duration float32
	var instrument string

	n, err := fmt.Sscan(strings.TrimSpace(line), &delay, &note, &duration, &instrument)
	if n != 4 || err != nil {
		fmt.Println(line)
		return synth.RawDelayedNoteData{}, errors.New("Failed to parse line data.")
	}

	return synth.RawDelayedNoteData{
		Delay:      delay,
		Note:       note,
		Duration:   duration,
		Instrument: instrument,
	}, nil
}

// Constructing a note arrangement from a file.
func constructNoteArrangement(file *os.File) ([]synth.RawDelayedNoteData, error) {
	na := []synth.RawDelayedNoteData{}
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if line == "\n" {
			continue
		} else if err == io.EOF {
			break
		} else if err != nil {
			return []synth.RawDelayedNoteData{}, errors.New("Encountered an error while reading a line: " + err.Error())
		}

		dnd, err := parseNoteData(line)
		if err != nil {
			return []synth.RawDelayedNoteData{}, errors.New("Failed to parse note data: " + err.Error())
		}

		na = append(na, dnd)
	}

	return na, nil
}

// Attempting to load in a file of the old format and write it out to the new
// location in the new format.
func ConvertFiles(from, to string) error {
	source, err := os.Open(from)
	if err != nil {
		return errors.New("Failed to open source file: " + err.Error())
	}
	defer source.Close()

	na, err := constructNoteArrangement(source)
	if err != nil {
		return errors.New("Failed to construct note arrangement: " + err.Error())
	}

	dest, err := os.Create(to)
	if err != nil {
		return errors.New("Failed to open destination file: " + err.Error())
	}
	defer dest.Close()

	enc := json.NewEncoder(dest)
	err = enc.Encode(na)
	if err != nil {
		return errors.New("Failed to encode JSON: " + err.Error())
	}

	return nil
}
