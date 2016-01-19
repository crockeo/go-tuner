package midi

// A MIDI chunk.
type Chunk struct {
	Title  string
	Length uint32
	Bytes  []byte
}

// The header data for a MIDI file.
type Header struct {
	Format   uint16
	Tracks   uint16
	Division int16
}

// A set of MIDI events that constitute a track.
type Track []Event

// A single MIDI event.
type Event struct {
	Delay int
}

// The entire structure of a MIDI file.
type MIDI struct {
	Header Header
	Tracks []Track
}
