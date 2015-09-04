package synth

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
