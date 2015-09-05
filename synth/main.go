package synth

import (
	"time"
)

// Starting a synth with a beginning note arrangement.
func StartSynthWith(na *NoteArrangement, noteChannel chan DelayedNoteData) error {
	var pd *PrimaryDriver
	if na == nil {
		pd = NewPrimaryDriverEmpty()
	} else {
		pd = NewPrimaryDriver(*na)
	}

	errChannel := make(chan error)
	defer close(errChannel)

	exitChannel := make(chan bool)
	defer close(exitChannel)

	go RunSynth(pd, errChannel, exitChannel)
	for {
		select {
		case err := <-errChannel:
			if err != nil {
				exitChannel <- true
				return err
			}
		}

		dnd := <-noteChannel
		pd.AddDelayedNote(dnd)

		time.Sleep(50 * time.Millisecond)
	}

	exitChannel <- true
	return nil
}

// Starting the synth with a channel for note data.
func StartSynth(noteChannel chan DelayedNoteData) error {
	return StartSynthWith(nil, noteChannel)
}
