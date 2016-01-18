package midi

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

// Shortcut for binary.Read.
func convertBytes(bs []byte, bo binary.ByteOrder, data interface{}) {
	buf := bytes.NewBuffer(bs)
	binary.Read(buf, bo, data)
}

// Reading a single byte.
func readByte(reader io.Reader) (byte, error) {
	b := make([]byte, 1)
	_, err := reader.Read(b)
	if err != nil {
		return 0x0, err
	}

	return b[0], nil
}

// Reading a variable-length string of bytes from a reader per its MIDI
// specification.
func readVarBytes(reader io.Reader) ([]byte, error) {
	bytes := []byte{}

	var err error
	for b, err := readByte(reader); err == nil; b, err = readByte(reader) {
		bytes = append(bytes, b)

		if b < 0x80 {
			break
		}
	}

	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

// Reading a chunk from a file.
func ReadChunk(reader io.Reader) (Chunk, error) {
	titleBytes := make([]byte, 4)
	_, err := reader.Read(titleBytes)
	if err != nil {
		return Chunk{}, err
	}

	lengthBytes := make([]byte, 4)
	_, err = reader.Read(lengthBytes)
	if err != nil {
		return Chunk{}, err
	}
	var length uint32
	convertBytes(lengthBytes, binary.BigEndian, &length)

	dataBytes := make([]byte, length)

	return Chunk{
		string(titleBytes),
		length,
		dataBytes,
	}, nil
}

// Reading the initial MIDI header from a file.
func ReadHeader(reader io.Reader) (Header, error) {
	chunk, err := ReadChunk(reader)
	if err != nil {
		return Header{}, err
	} else if chunk.Title != "MThd" || chunk.Length != 6 {
		return Header{}, errors.New("Invalid chunk header.")
	}

	var format, tracks uint16
	var division int16

	convertBytes(chunk.Bytes[0:2], binary.BigEndian, &format)
	convertBytes(chunk.Bytes[2:4], binary.BigEndian, &tracks)
	convertBytes(chunk.Bytes[4:6], binary.BigEndian, &division)

	return Header{
		format,
		tracks,
		division,
	}, nil
}

// Reading an event from in from a reader.
func ReadEvent(reader io.Reader) (Event, error) {
	// TODO:
	//  1. Read in the delta_time.
	//  2. Read in the event data.
	//  3. Compile that into an event.

	delayBytes, err := readVarBytes(reader)
	if err != nil {
		return nil, err
	}
	var delay int
	convertBytes(delayBytes, binary.BigEndian, &delay)

	b, err := readByte(reader)
	if err != nil {
		return nil, err
	}

	var e Event
	switch b {
	// Loading a sysex event.
	case 0xF0, 0xF7:
		lenBytes, err := readVarBytes(reader)
		var len int
		convertBytes(lenBytes, binary.BigEndian, &len)

		data := make([]byte, len)
		_, err = reader.Read(data)
		if err != nil {
			return nil, err
		}

		return SysexEvent{
			delay,
			b == 0xF7,
			data,
		}, nil
	// Loading a meta event.
	case 0xFF:
		t, err := readByte(reader)
		if err != nil {
			return nil, err
		}

		switch t {
		case 0x00:
		case 0x01:
		case 0x02:
		case 0x03:
		case 0x04:
		case 0x05:
		case 0x06:
		case 0x07:
		case 0x20:
		case 0x2F:
		case 0x51:
		case 0x54:
		case 0x58:
		case 0x59:
		case 0x7F:
		}
	default:
		e, err = ReadMIDIEvent(reader, b)
	}

	if err != nil {
		return nil, err
	}

	return e, nil
}

// Reading a given track in from a reader.
func ReadTrack(reader io.Reader) (Track, error) {
	chunk, err := ReadChunk(reader)
	if err != nil {
		return Track{}, err
	} else if chunk.Title != "MTrk" {
		return Track{}, errors.New("Invalid chunk header.")
	}

	track := Track{}
	buf := bytes.NewBuffer(chunk.Bytes)
	for buf.Len() > 0 {
		event, err := ReadEvent(buf)
		if err != nil {
			return Track{}, nil
		}

		// Event Kind != a MIDI sound event, either a sysex event or a meta
		// event. Might re-add later, but it simplifies things.
		if event.Kind() == 1 || event.Kind() == 2 {
			continue
		}

		track = append(track, event)
	}

	return Track{}, nil
}

// Constructing and returning a struct from the data contained in an io.Reader.
func Read(reader io.Reader) (*MIDI, error) {
	header, err := ReadHeader(reader)
	if err != nil {
		return nil, err
	}

	tracks := make([]Track, header.Tracks)
	for i := 0; i < int(header.Tracks); i++ {
		track, err := ReadTrack(reader)
		if err != nil {
			return nil, err
		}

		tracks[i] = track
	}

	return nil, errors.New("Read not yet implemented")
}

// Constructing and returning a struct from a file on disk.
func ReadFile(path string) (*MIDI, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Read(file)
}

// Writing a MIDI structure out to some io.Writer.
func (m *MIDI) Write(writer io.Writer) error {
	return errors.New("Write not yet implemented.")
}

// Writing a MIDI structure out to a file location on disk. Opens a file and
// calls Write(..)
func (m *MIDI) WriteFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return m.WriteFile(path)
}
