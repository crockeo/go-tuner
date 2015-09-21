package synth

import (
	"time"
)

// The function to start a synth with the intent of being asynchronous WITH a
// given slice of starting notes.
//
// na           - The slice of notes to play.
// iNoteChannel - A channel to provide note data.
// iQuitChannel - A channel to query for an external exit signal.
// oErrChannel  - A channel to send out error information to the calling
//                function.
func StartSynthAsyncWith(na *NoteArrangement, iNoteChannel chan DelayedNoteData, ioQuitChannel chan bool, oErrChannel chan error, quitWhenDone bool) {
	var pd *PrimaryDriver
	if na == nil {
		pd = NewPrimaryDriverEmpty()
	} else {
		pd = NewPrimaryDriver(*na)
	}

	exitChannel := make(chan bool)
	defer close(exitChannel)

	errChannel := make(chan error)
	defer close(errChannel)

	go RunSynth(pd, errChannel, quitWhenDone, exitChannel)

	// Using a label here so we can break out of the for loop from inside the
	// case statement.
	//
	// The case statement is used so we can aggressively scan for information
	// from all three channels.
	for {
		select {
		case dnd := <-iNoteChannel:
			pd.AddDelayedNote(dnd)
		case _ = <-ioQuitChannel:
			return
		case _ = <-exitChannel:
			ioQuitChannel <- true
			return
		case err := <-errChannel:
			if err != nil {
				exitChannel <- true
				oErrChannel <- err
			}
		}

		time.Sleep(50 * time.Millisecond)
	}
}

// The function to start a synth with the intent of being asynchronous.
//
// iNoteChannel - A channel to provide note data.
// iQuitChannel - A channel to query for an external exit signal.
// oErrChannel  - A channel to send out error information to the calling
//                function.
func StartSynthAsync(iNoteChannel chan DelayedNoteData, iQuitChannel chan bool, oErrChannel chan error) {
	StartSynthAsyncWith(nil, iNoteChannel, iQuitChannel, oErrChannel, false)
}

// Starting a synth with a beginning note arrangement.
func StartSynthWith(na *NoteArrangement, iNoteChannel chan DelayedNoteData, quitWhenDone bool) error {
	iQuitChannel := make(chan bool)
	defer close(iQuitChannel)

	errChannel := make(chan error)
	defer close(errChannel)

	go StartSynthAsyncWith(na, iNoteChannel, iQuitChannel, errChannel, quitWhenDone)

	for {
		select {
		case _ = <-iQuitChannel:
			return nil
		case err := <-errChannel:
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Starting the synth with a channel for note data.
func StartSynth(noteChannel chan DelayedNoteData) error {
	return StartSynthWith(nil, noteChannel, false)
}
