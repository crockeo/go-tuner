package synth

import (
	"bytes"
	"encoding/json"
	"errors"
)

// Type RawDelayedNoteData is raw data from a message that can be converted into
// DelayedNoteData after decoding its JSON.
type RawDelayedNoteData struct {
	Delay      float64 `json:"delay"`
	Note       string  `json:"note"`
	Duration   float64 `json:"duration"`
	Instrument string  `json:"instrument"`
}

//// Attempting to parse a given message into a DelayedNote. It will return
//// an error upon failure.
func ParseMessage(str string) (DelayedNoteData, error) {
	var rdnd RawDelayedNoteData
	dec := json.NewDecoder(bytes.NewBufferString(str))
	err := dec.Decode(&rdnd)
	if err != nil {
		return DelayedNoteData{}, errors.New("Could not parse JSON information: " + err.Error())
	}

	// TODO: Look up instrument and note from some map.

	return DelayedNoteData{}, nil
}
