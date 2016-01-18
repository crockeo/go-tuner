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
type Event interface {
	Delay() int
	Kind() int
}

// A MIDI event.
type MIDIEvent struct {
	Ticks int
}

func (me MIDIEvent) Delay() int { return me.Ticks }
func (me MIDIEvent) Kind() int  { return 0 }

// A sysex event.
type SysexEvent struct {
	Ticks  int
	Escape bool
	Data   []byte
}

func (se SysexEvent) Delay() int { return se.Ticks }
func (se SysexEvent) Kind() int  { return 1 }

// A meta event.
type MetaEvent struct {
	Ticks int
	Type  int
}

func (me MetaEvent) Delay() int { return me.Ticks }
func (me MetaEvent) Kind() int  { return 2 }

// The entire structure of a MIDI file.
type MIDI struct {
	Header Header
	Tracks []Track
}
