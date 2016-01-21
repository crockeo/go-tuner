package filestore

import (
	"errors"
	"fmt"
	"github.com/crockeo/go-tuner/filestore/midi"
	"github.com/crockeo/go-tuner/synth"
	"io"
)

// Converting between a key (from MIDI input) and the corresponding string that
// is used elsewhere in the program.
func numToNote(key uint8) string {
	octave := key / 12
	note := key - (octave * 12)

	var str string
	switch note {
	case 0:
		str = "C"
	case 1:
		str = "C#"
	case 2:
		str = "D"
	case 3:
		str = "D#"
	case 4:
		str = "E"
	case 5:
		str = "F"
	case 6:
		str = "F#"
	case 7:
		str = "G"
	case 8:
		str = "G#"
	case 9:
		str = "A"
	case 10:
		str = "A#"
	case 11:
		str = "B"
	}

	return fmt.Sprintf("%s%d", str, octave)
}

// Converting a tick amount to a real time delay.
func convertTick(division int16, delay uint) float32 {
	switch division >> 15 {
	case 0:
		return 0.5 * (float32(delay) / float32(division&0x7FFF))
	case 1:
	default:
		panic("convertTick defualt select - you shouldn't be able to get here.")
	}

	return 0.0
}

// Constructing a single []synth.RawDelayedNoteData from a MIDI track.
func constructNoteArrangement(header midi.Header, track midi.Track) []synth.RawDelayedNoteData {
	rdnds := []synth.RawDelayedNoteData{}

	for i := 0; i < len(track); i++ {
		t1 := track[i]

		var diff uint = 0
		for j := i + 1; j < len(track); j++ {
			diff += track[j].Delay
			if !track[j].Switch && track[j].Channel == t1.Channel && track[j].Key == t1.Key {
				break
			}
		}

		// TODO: Non-fixed note duration.

		delay := convertTick(header.Division, t1.Delay)
		rdnds = append(rdnds, synth.RawDelayedNoteData{
			delay,
			numToNote(t1.Key),
			0.7, //(convertTick(header.Division, diff) - delay) + 0.25,
			"guitar",
		})
	}

	return rdnds
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

			continue
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

			continue
		}

		// Appending whichever note is going to happen soonest.
		if side == 0 {
			if na1[0].Delay < na2[0].Delay {
				accum = append(accum, na1[0])
				offset = na1[0].Delay
				na1 = na1[1:]
				side = 1
			} else {
				accum = append(accum, na2[0])
				offset = na2[0].Delay
				na2 = na2[1:]
				side = 2
			}
		} else if side == 1 {
			if na1[0].Delay < na2[0].Delay-offset {
				accum = append(accum, na1[0])
				offset += na1[0].Delay
				na1 = na1[1:]
			} else {
				accum = append(accum, na2[0])
				offset = na2[0].Delay
				na2 = na2[1:]
				side = 2
			}
		} else if side == 2 {
			if na1[0].Delay-offset < na2[0].Delay {
				accum = append(accum, na1[0])
				offset = na1[0].Delay
				na1 = na1[1:]
				side = 1
			} else {
				accum = append(accum, na2[0])
				offset += na2[0].Delay
				na2 = na2[1:]
			}
		}
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
