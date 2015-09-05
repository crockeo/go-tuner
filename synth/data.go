package synth

import (
	"bytes"
	"encoding/json"
	"errors"
)

var (
	// The map from names of notes to their frequencies.
	notes = map[string]float64{
		"C0":  16.35,
		"C#0": 17.32,
		"Db0": 17.32,
		"D0":  18.35,
		"D#0": 19.45,
		"Eb0": 19.45,
		"E0":  20.60,
		"F0":  21.83,
		"F#0": 23.12,
		"Gb0": 23.12,
		"G0":  24.50,
		"G#0": 25.96,
		"Ab0": 25.96,
		"A0":  27.50,
		"A#0": 29.14,
		"Bb0": 29.14,
		"B0":  30.87,
		"C1":  32.70,
		"C#1": 34.65,
		"Db1": 34.65,
		"D1":  36.71,
		"D#1": 38.89,
		"Eb1": 38.89,
		"E1":  41.20,
		"F1":  43.65,
		"F#1": 46.25,
		"Gb1": 46.25,
		"G1":  49.00,
		"G#1": 51.91,
		"Ab1": 51.91,
		"A1":  55.00,
		"A#1": 58.27,
		"Bb1": 58.27,
		"B1":  61.74,
		"C2":  65.41,
		"C#2": 69.30,
		"Db2": 69.30,
		"D2":  73.42,
		"D#2": 77.78,
		"Eb2": 77.78,
		"E2":  82.41,
		"F2":  87.31,
		"F#2": 92.50,
		"Gb2": 92.50,
		"G2":  98.00,
		"G#2": 103.83,
		"Ab2": 103.83,
		"A2":  110.00,
		"A#2": 116.54,
		"Bb2": 116.54,
		"B2":  123.47,
		"C3":  130.81,
		"C#3": 138.59,
		"Db3": 138.59,
		"D3":  146.83,
		"D#3": 155.56,
		"Eb3": 155.56,
		"E3":  164.81,
		"F3":  174.61,
		"F#3": 185.00,
		"Gb3": 185.00,
		"G3":  196.00,
		"G#3": 207.65,
		"Ab3": 207.65,
		"A3":  220.00,
		"A#3": 233.08,
		"Bb3": 233.08,
		"B3":  246.94,
		"C4":  261.63,
		"C#4": 277.18,
		"Db4": 277.18,
		"D4":  293.66,
		"D#4": 311.13,
		"Eb4": 311.13,
		"E4":  329.63,
		"F4":  349.23,
		"F#4": 369.99,
		"Gb4": 369.99,
		"G4":  392.00,
		"G#4": 415.30,
		"Ab4": 415.30,
		"A4":  440.00,
		"A#4": 466.16,
		"Bb4": 466.16,
		"B4":  493.88,
		"C5":  523.25,
		"C#5": 554.37,
		"Db5": 554.37,
		"D5":  587.33,
		"D#5": 622.25,
		"Eb5": 622.25,
		"E5":  659.25,
		"F5":  698.46,
		"F#5": 739.99,
		"Gb5": 739.99,
		"G5":  783.99,
		"G#5": 830.61,
		"Ab5": 830.61,
		"A5":  880.00,
		"A#5": 932.33,
		"Bb5": 932.33,
		"B5":  987.77,
		"C6":  1046.50,
		"C#6": 1108.73,
		"Db6": 1108.73,
		"D6":  1174.66,
		"D#6": 1244.51,
		"Eb6": 1244.51,
		"E6":  1318.51,
		"F6":  1396.91,
		"F#6": 1479.98,
		"Gb6": 1479.98,
		"G6":  1567.98,
		"G#6": 1661.22,
		"Ab6": 1661.22,
		"A6":  1760.00,
		"A#6": 1864.66,
		"Bb6": 1864.66,
		"B6":  1975.53,
		"C7":  2093.00,
		"C#7": 2217.46,
		"Db7": 2217.46,
		"D7":  2349.32,
		"D#7": 2489.02,
		"Eb7": 2489.02,
		"E7":  2637.02,
		"F7":  2793.83,
		"F#7": 2959.96,
		"Gb7": 2959.96,
		"G7":  3135.96,
		"G#7": 3322.44,
		"Ab7": 3322.44,
		"A7":  3520.00,
		"A#7": 3729.31,
		"Bb7": 3729.31,
		"B7":  3951.07,
		"C8":  4186.01,
		"C#8": 4434.92,
		"Db8": 4434.92,
		"D8":  4698.63,
		"D#8": 4978.03,
		"Eb8": 4978.03,
		"E8":  5274.04,
		"F8":  5587.65,
		"F#8": 5919.91,
		"Gb8": 5919.91,
		"G8":  6271.93,
		"G#8": 6644.88,
		"Ab8": 6644.88,
		"A8":  7040.00,
		"A#8": 7458.62,
		"Bb8": 7458.62,
		"B8":  7902.13,
	}

	// The map of names of instruments to their NoteData-generating functions.
	instruments = map[string]func(float64, float64, float64) NoteData{
		"guitar": GuitarNote,
	}
)

// Type Overtone represents a relationship to a primary not that is some
// multiplier on the note's frequency with a given volume.
type Overtone struct {
	Relation float64
	Volume   float64
}

// Type NoteData represents the information necessary to play a single note from
// a given instrument.
type NoteData struct {
	Duration  float64
	Volume    float64
	Frequency float64
	FadeFunc  func(float64, float64) float64
	Overtones []Overtone
}

// Creating a NoteData from a guitar note.
func GuitarNote(duration, volume, frequency float64) NoteData {
	return NoteData{
		Duration:  duration,
		Volume:    volume,
		Frequency: frequency,

		FadeFunc: func(time, duration float64) float64 {
			return 1.0 - (time / duration)
		},

		Overtones: []Overtone{
			Overtone{2, 0.500},
			Overtone{3, 0.333},
			Overtone{4, 0.250},
			Overtone{5, 0.200},
			Overtone{6, 0.166},
			Overtone{7, 0.142},
			Overtone{8, 0.125},
		},
	}
}

// Type rawDelayedNoteData is raw data from a message that can be converted into
// DelayedNoteData after decoding its JSON.
type rawDelayedNoteData struct {
	Delay      float64 `json:"delay"`
	Note       string  `json:"note"`
	Duration   float64 `json:"duration"`
	Instrument string  `json:"instrument"`
}

// Type DelayedNoteData is a container that houses the delay and the note data
// for a given note.
type DelayedNoteData struct {
	Delay float64
	ND    NoteData
}

// Converting a []byte representation of a DelayedNoteData into a structured
// piece of data.
func (dnd DelayedNoteData) UnmarshalJSON(data []byte) error {
	var rdnd rawDelayedNoteData

	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&rdnd)
	if err != nil {
		return err
	}

	note, ok := notes[rdnd.Note]
	if !ok {
		return errors.New("Invalid note name.")
	}

	instrument, ok := instruments[rdnd.Instrument]
	if !ok {
		return errors.New("Invalid instrument name.")
	}

	dnd.Delay = rdnd.Delay
	dnd.ND = instrument(rdnd.Duration, 1.0, note)

	return nil
}

// Type NoteArrangement is a synonym for an array of structs, one being the
// delay from the last note, the other being the associaed NoteData for this
// note.
type NoteArrangement []DelayedNoteData

// Creating an empty NoteArrangement.
func EmptyNoteArrangement() *NoteArrangement {
	na := new(NoteArrangement)
	return na
}

// Adding a piece of NoteData to a NoteArrangement at a certain delay.
func (na *NoteArrangement) AddNoteData(delay float64, nd NoteData) {
	_ = append(*na, DelayedNoteData{
		Delay: delay,
		ND:    nd,
	})
}
