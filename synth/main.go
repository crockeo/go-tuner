package synth

import (
	"time"
)

// Starting a synth with a beginning note arrangement.
func StartSynthWith(na *NoteArrangement, noteChannel chan DelayedNoteData, quitWhenDone bool) error {
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

	go RunSynth(pd, errChannel, quitWhenDone, exitChannel)

	// Using a label here so we can break out of the for loop from inside the
	// case statement.
	//
	// The case statement is used so we can aggressively scan for information
	// from all three channels.
outer:
	for {
		select {
		case err := <-errChannel:
			if err != nil {
				exitChannel <- true
				return err
			}
		case dnd := <-noteChannel:
			pd.AddDelayedNote(dnd)
		case _ = <-exitChannel:
			break outer
		}

		time.Sleep(50 * time.Millisecond)
	}

	return nil
}

// Starting the synth with a channel for note data.
func StartSynth(noteChannel chan DelayedNoteData) error {
	return StartSynthWith(nil, noteChannel, false)
}
