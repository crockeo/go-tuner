package midi

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
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

// Reading in a variable quantity int.
func varInt(reader io.Reader) (int, error) {
	vbs, err := readVarBytes(reader)
	if err != nil {
		return 0, err
	}

	var n int
	convertBytes(vbs, binary.BigEndian, &n)

	return n, nil
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
	_, err = reader.Read(dataBytes)
	if err != nil {
		return Chunk{}, err
	}

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
func ReadEvent(reader io.Reader) (Event, bool, error) {
	delayBytes, err := readVarBytes(reader)
	if err != nil {
		return Event{}, false, err
	}
	var delay int
	convertBytes(delayBytes, binary.BigEndian, &delay)

	b, err := readByte(reader)
	if err != nil {
		return Event{}, false, err
	}

	var e Event
	switch b {
	// Loading a sysex event.
	case 0xF0, 0xF7:
		len, err := varInt(reader)
		if err != nil {
			return Event{}, false, err
		}

		data := make([]byte, len)
		_, err = reader.Read(data)
		if err != nil {
			return Event{}, false, err
		}

		return Event{}, true, nil
	// Loading a meta event.
	case 0xFF:
		t, err := readByte(reader)
		if err != nil {
			return Event{}, false, err
		}

		switch t {
		case 0x00:
			bs := make([]byte, 3)
			_, err := reader.Read(bs)
			if err != nil {
				return Event{}, false, err
			}
		case 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x7F:
			len, err := varInt(reader)
			if err != nil {
				return Event{}, false, err
			}

			bs := make([]byte, len)
			if _, err := reader.Read(bs); err != nil {
				return Event{}, false, err
			}

			return Event{}, true, nil
		case 0x20:
			bs := make([]byte, 2)
			if _, err := reader.Read(bs); err != nil {
				return Event{}, false, err
			}

			return Event{}, true, nil
		case 0x2F:
			bs := make([]byte, 1)
			if _, err := reader.Read(bs); err != nil {
				return Event{}, false, err
			}

			return Event{}, true, nil
		case 0x51:
			bs := make([]byte, 4)
			if _, err := reader.Read(bs); err != nil {
				return Event{}, false, err
			}

			return Event{}, true, nil
		case 0x54:
			bs := make([]byte, 6)
			if _, err := reader.Read(bs); err != nil {
				return Event{}, false, err
			}

			return Event{}, true, nil
		case 0x58:
			bs := make([]byte, 5)
			if _, err := reader.Read(bs); err != nil {
				return Event{}, false, err
			}

			return Event{}, true, nil
		case 0x59:
			bs := make([]byte, 3)
			if _, err := reader.Read(bs); err != nil {
				return Event{}, false, err
			}

			return Event{}, true, nil
		}
	default:
		kind := b >> 4

		switch kind {
		case 0x8, 0x9:
			s := kind == 0x9

			key, err := readByte(reader)
			if err != nil {
				return Event{}, false, err
			}

			velocity, err := readByte(reader)
			if err != nil {
				return Event{}, false, err
			}

			if s && velocity == 0 {
				s = false
				velocity = 44
			}

			return Event{
				delay,
				s,
				uint8(b & 0x0F),
				uint8(key),
				uint8(velocity),
			}, false, nil
		case 0xA, 0xB, 0xE:
			bs := make([]byte, 2)
			if _, err := reader.Read(bs); err != nil {
				return Event{}, false, err
			}

			return Event{}, true, nil
		case 0xC, 0xD:
			bs := make([]byte, 1)
			if _, err := reader.Read(bs); err != nil {
				return Event{}, false, err
			}

			return Event{}, true, nil
		}
	}

	if err != nil {
		return Event{}, false, err
	}

	return e, false, nil
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
		event, skip, err := ReadEvent(buf)
		if err != nil {
			return Track{}, nil
		}

		if skip {
			continue
		}

		track = append(track, event)
	}

	return track, nil
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

	return &MIDI{
		header,
		tracks,
	}, nil
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
	fmt.Println("Writing MIDI files is not yet supported.")
	return nil
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
