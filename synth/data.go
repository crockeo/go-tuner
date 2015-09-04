package synth

const (
	// A conversion between note names and note frequencies.
	C0  float64 = 16.35
	Db0 float64 = 17.32
	D0  float64 = 18.35
	Eb0 float64 = 19.45
	E0  float64 = 20.60
	F0  float64 = 21.83
	Gb0 float64 = 23.12
	G0  float64 = 24.50
	Ab0 float64 = 25.96
	A0  float64 = 27.50
	Bb0 float64 = 29.14
	B0  float64 = 30.87
	C1  float64 = 32.70
	Db1 float64 = 34.65
	D1  float64 = 36.71
	Eb1 float64 = 38.89
	E1  float64 = 41.20
	F1  float64 = 43.65
	Gb1 float64 = 46.25
	G1  float64 = 49.00
	Ab1 float64 = 51.91
	A1  float64 = 55.00
	Bb1 float64 = 58.27
	B1  float64 = 61.74
	C2  float64 = 65.41
	Db2 float64 = 69.30
	D2  float64 = 73.42
	Eb2 float64 = 77.78
	E2  float64 = 82.41
	F2  float64 = 87.31
	Gb2 float64 = 92.50
	G2  float64 = 98.00
	Ab2 float64 = 103.83
	A2  float64 = 110.00
	Bb2 float64 = 116.54
	B2  float64 = 123.47
	C3  float64 = 130.81
	Db3 float64 = 138.59
	D3  float64 = 146.83
	Eb3 float64 = 155.56
	E3  float64 = 164.81
	F3  float64 = 174.61
	Gb3 float64 = 185.00
	G3  float64 = 196.00
	Ab3 float64 = 207.65
	A3  float64 = 220.00
	Bb3 float64 = 233.08
	B3  float64 = 246.94
	C4  float64 = 261.63
	Db4 float64 = 277.18
	D4  float64 = 293.66
	Eb4 float64 = 311.13
	E4  float64 = 329.63
	F4  float64 = 349.23
	Gb4 float64 = 369.99
	G4  float64 = 392.00
	Ab4 float64 = 415.30
	A4  float64 = 440.00
	Bb4 float64 = 466.16
	B4  float64 = 493.88
	C5  float64 = 523.25
	Db5 float64 = 554.37
	D5  float64 = 587.33
	Eb5 float64 = 622.25
	E5  float64 = 659.25
	F5  float64 = 698.46
	Gb5 float64 = 739.99
	G5  float64 = 783.99
	Ab5 float64 = 830.61
	A5  float64 = 880.00
	Bb5 float64 = 932.33
	B5  float64 = 987.77
	C6  float64 = 1046.50
	Db6 float64 = 1108.73
	D6  float64 = 1174.66
	Eb6 float64 = 1244.51
	E6  float64 = 1318.51
	F6  float64 = 1396.91
	Gb6 float64 = 1479.98
	G6  float64 = 1567.98
	Ab6 float64 = 1661.22
	A6  float64 = 1760.00
	Bb6 float64 = 1864.66
	B6  float64 = 1975.53
	C7  float64 = 2093.00
	Db7 float64 = 2217.46
	D7  float64 = 2349.32
	Eb7 float64 = 2489.02
	E7  float64 = 2637.02
	F7  float64 = 2793.83
	Gb7 float64 = 2959.96
	G7  float64 = 3135.96
	Ab7 float64 = 3322.44
	A7  float64 = 3520.00
	Bb7 float64 = 3729.31
	B7  float64 = 3951.07
	C8  float64 = 4186.01
	Db8 float64 = 4434.92
	D8  float64 = 4698.63
	Eb8 float64 = 4978.03
	E8  float64 = 5274.04
	F8  float64 = 5587.65
	Gb8 float64 = 5919.91
	G8  float64 = 6271.93
	Ab8 float64 = 6644.88
	A8  float64 = 7040.00
	Bb8 float64 = 7458.62
	B8  float64 = 7902.13
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

// Type DelayedNoteData is a container that houses the delay and the note data
// for a given note.
type DelayedNoteData struct {
	Delay float64
	ND    NoteData
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
