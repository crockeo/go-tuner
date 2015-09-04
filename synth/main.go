package synth

import (
	"errors"
	"time"
)

// Starting the synth with static synth data.
func StartStaticSynth(na *NoteArrangement) error {
	if na == nil {
		return errors.New("Cannot work with a nil NoteArrangement.")
	}
	pd := NewPrimaryDriver(*na)

	return RunSynth(pd)
}

// Starting the synth with a channel for note data.
func StartDynamicSynth(noteChannel chan DelayedNoteData) error {
	pd := NewPrimaryDriverEmpty()

	errChannel := make(chan error)
	defer close(errChannel)

	exitChannel := make(chan bool)
	defer close(exitChannel)

	go RunSynthAsync(pd, errChannel, exitChannel)
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
