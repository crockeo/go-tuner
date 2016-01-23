package synth

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

// The base note to be used for the other calculations.
const A4 float32 = 440
const A4I int = 57

// Converting the int-representation of a note to the string-representation of
// a note.
func NoteToString(note int) string {
	octave := note / 12

	var str string
	switch note - (octave * 12) {
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

// Converting the string-representation of a note to the int-representation of a
// note.
func NoteToInt(note string) (int, error) {
	if len(note) < 2 {
		return 0, errors.New("Malformed note: " + note)
	}

	var edge int
	for edge = 0; edge < len(note); edge++ {
		if '0' <= note[edge] && note[edge] <= '9' {
			break
		}
	}

	oct, err := strconv.ParseInt(note[edge:], 10, 0)
	if err != nil {
		return 0, err
	}

	var n int
	switch note[:edge] {
	case "C":
		n = 0
	case "C#":
		n = 1
	case "D":
		n = 2
	case "D#":
		n = 3
	case "E":
		n = 4
	case "F":
		n = 5
	case "F#":
		n = 6
	case "G":
		n = 7
	case "G#":
		n = 8
	case "A":
		n = 9
	case "A#":
		n = 10
	case "B":
		n = 11
	default:
		return 0, errors.New("Invalid note name: " + note[:edge])
	}

	return int(oct)*12 + n, nil
}

// Calculating the frequency of a given note.
func CalculateFrequency(note int) float32 {
	return A4 * float32(math.Pow(1.059463094359, float64(note-A4I)))
}

// Calculating the frequency of a given note from its string representation.
func CalculateFrequencyStr(note string) (float32, error) {
	n, err := NoteToInt(note)
	if err != nil {
		return 0.0, err
	}

	return CalculateFrequency(n), nil
}
