package synth

import "errors"

// Starting the synth with static synth data.
func StartStaticSynth(na *NoteArrangement) error {
	return errors.New("StartStaticSynth not yet implemented.")
}

// Starting the synth with a channel for note data.
func StartDynamicSynth(noteChannel chan DelayedNoteData) error {
	return errors.New("StartDynamicSynth not yet implemented.")
}
