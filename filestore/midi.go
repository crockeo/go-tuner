package filestore

import (
	"errors"
	"github.com/crockeo/go-tuner/filestore/midi"
	"github.com/crockeo/go-tuner/synth"
	"io"
)

// Constructing a single []synth.RawDelayedNoteData from a MIDI track.
func constructNoteArrangement(header midi.Header, track midi.Track) []synth.RawDelayedNoteData {
	// TODO: Construct a []synth.RawDelayedNoteData from a track.
	return []synth.RawDelayedNoteData{}
}

// Merging two []synth.RawDelayedNoteDatas together.
func mergeNoteArrangements(na1 []synth.RawDelayedNoteData, na2 []synth.RawDelayedNoteData) []synth.RawDelayedNoteData {
	accum := []synth.RawDelayedNoteData{}
	offset := float32(0.0)
	side := 0

	for len(na1) > 0 || len(na2) > 0 {
		// There are no more notes from na1 to append, so appending everything
		// that is left from na2.
		if len(na1) == 0 {
			t := na2[0]
			if side == 1 {
				t.Delay -= offset

				offset = 0.0
				side = 0
			}

			accum = append(accum, t)
			na2 = na2[1:]
		}

		// There are no more notes from na2 to append, so appending everything
		// that is left from na1.
		if len(na2) == 0 {
			t := na1[0]
			if side == 2 {
				t.Delay -= offset

				offset = 0.0
				side = 0
			}

			accum = append(accum, t)
			na2 = na2[1:]
		}

		// Appending whichever note is going to happen soonest.
	}

	return accum
}

// Dealing with synth.RawDelayedNoteData from a MIDI file.
type MIDIArrangement struct{}

func (a MIDIArrangement) ReadNoteArrangement(reader io.Reader) ([]synth.RawDelayedNoteData, error) {
	m, err := midi.Read(reader)
	if err != nil {
		return []synth.RawDelayedNoteData{}, err
	}

	switch m.Header.Format {
	case 0:
		if len(m.Tracks) != 1 {
			return []synth.RawDelayedNoteData{}, errors.New("Malformed track length.")
		}

		return constructNoteArrangement(m.Header, m.Tracks[0]), nil
	case 1:
		accum := []synth.RawDelayedNoteData{}
		for i := 0; i < len(m.Tracks); i++ {
			accum = mergeNoteArrangements(accum, constructNoteArrangement(m.Header, m.Tracks[i]))
		}

		return accum, nil
	case 2:
		accum := []synth.RawDelayedNoteData{}
		for i := 0; i < len(m.Tracks); i++ {
			accum = append(accum, constructNoteArrangement(m.Header, m.Tracks[i])...)
		}

		return accum, nil
	default:
		return []synth.RawDelayedNoteData{}, errors.New("Unexpected format.")
	}

	return []synth.RawDelayedNoteData{}, nil
}

func (a MIDIArrangement) WriteNoteArrangement(writer io.Writer, notes []synth.RawDelayedNoteData) error {
	return errors.New("WriteNoteArrangement not yet implemented.")
}
