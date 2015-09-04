package synth

import (
	"github.com/HardWareGuy/portaudio-go"
	"math"
)

// Type Driver is an interface to define the required behavior for a data
// structure that can be used in driving PortAudio music synthesis.
type Driver interface {
	// Getting the number of output channels this driver is expecting.
	OutputChannels() int

	// Calculating the output on whatever set of channels for a given driver.
	CalculateOutput() []float64

	// Stepping the internal phases given a sample rate.
	StepPhases(int)
}

// The driver to play a single note for a given duration.
type SingleDriver struct {
	Note      NoteData // The data for this note.
	Phase     float64  // The current phase of the driver.
	Time      float64  // The current time of the SingleDriver.
	StartTime float64  // The time this driver was started - if used without a PrimaryDriver it will always be 0.
}

// Creating a new SingleDriver from a given NoteData.
func NewSingleDriver(nd NoteData) *SingleDriver {
	sd := new(SingleDriver)

	sd.Note = nd
	sd.Phase = 0
	sd.Time = 0
	sd.StartTime = 0

	return sd
}

// Creating a new SingleDriver to be used as a player inside of a PrimaryDriver
// from a given NoteData.
func NewSingleDriverChild(nd NoteData, startTime float64) *SingleDriver {
	sd := new(SingleDriver)

	sd.Note = nd
	sd.Phase = 0
	sd.Time = startTime
	sd.StartTime = startTime

	return sd
}

// Getting the number of output channels this driver is expecting.
func (sd *SingleDriver) OutputChannels() int {
	return 2
}

// Calculating the output on whatever set of channels for a given driver.
func (sd *SingleDriver) CalculateOutput() []float64 {
	output := math.Sin(sd.Phase) * sd.Note.Volume * sd.Note.FadeFunc(sd.Time-sd.StartTime, sd.Note.Duration)

	outputs := make([]float64, sd.OutputChannels())
	for k, _ := range outputs {
		outputs[k] = output
	}

	return outputs
}

// Stepping the internal phases given a sample rate.
func (sd *SingleDriver) StepPhases(sampleRate int) {
	sd.Phase += sd.Note.Frequency * (1 / float64(sampleRate))
	if sd.Phase >= 2*math.Pi {
		sd.Phase -= 2 * math.Pi
	}

	sd.Time += 1.0 / float64(sampleRate)
}

// The primary driver that is used by the rest of the program by default to
// start whichever synth.
type PrimaryDriver struct {
	QueuedNotes  []DelayedNoteData // The list of NoteDatas to add.
	CurrentNotes []*SingleDriver   // The list of current SingleDrivers.
	Time         float64           // The current time of the PrimaryDriver.
	LastTime     float64           // The time that the last SingleDriver was added.
}

// Creating a PrimaryDriver from existent data.
func NewPrimaryDriver(queuedNotes []DelayedNoteData) *PrimaryDriver {
	pd := new(PrimaryDriver)

	pd.QueuedNotes = queuedNotes
	pd.CurrentNotes = []*SingleDriver{}
	pd.Time = 0.0
	pd.LastTime = 0.0

	return pd
}

// Creating a PrimaryDriver with no information inside yet.
func NewPrimaryDriverEmpty() *PrimaryDriver {
	pd := new(PrimaryDriver)

	pd.QueuedNotes = []DelayedNoteData{}
	pd.CurrentNotes = []*SingleDriver{}
	pd.Time = 0.0
	pd.LastTime = 0.0

	return pd
}

// Trying to add a new DelayedNoteData to the list of queued notes inside of a
// PrimaryDriver.
func (pd *PrimaryDriver) AddDelayedNote(dnd DelayedNoteData) {

}

// Getting the number of output channels this driver is expecting.
func (pd *PrimaryDriver) OutputChannels() int { return 2 }

// Calculating the output on whatever set of channels for a given driver.
func (pd *PrimaryDriver) CalculateOutput() []float64 {
	v := 0.0

	return []float64{v, v}
}

// Stepping the internal phases given a sample rate.
func (pd *PrimaryDriver) StepPhases(sampleRate int) {
	// Appending new notes to the set of current notes.
	for pd.Time-pd.LastTime <= pd.QueuedNotes[0].Delay {
		pd.CurrentNotes = append(pd.CurrentNotes, NewSingleDriverChild(pd.QueuedNotes[0].ND, pd.Time))
		pd.QueuedNotes = pd.QueuedNotes[1:]

		pd.LastTime = pd.Time
	}

	// TODO: Remove elements from CurrentNotes that are terminated.

	// Stepping the phases for the sub drivers.
	for _, sd := range pd.CurrentNotes {
		sd.StepPhases(sampleRate)
	}

	pd.Time += 1.0 / float64(sampleRate)
}

// Returning a function to drive music synthesis given a driver and a sample
// rate.
func DriverFunction(driver Driver, sampleRate int) func([][]float64) {
	return func(out [][]float64) {
		for i := range out[0] {
			output := driver.CalculateOutput()
			for j := range output {
				out[j][i] = output[j]
			}

			driver.StepPhases(sampleRate)
		}
	}
}

// Asynchronously running a synth.
func RunSynthAsync(driver Driver, errChannel chan error, exitChannel chan bool) {
	portaudio.Initialize()
	defer portaudio.Terminate()

	// TODO: Rest of this stuff.

	<-exitChannel
}

// Running a synth from a given Driver.
func RunSynth(driver Driver) error {
	portaudio.Initialize()
	defer portaudio.Terminate()

	// TODO: Rest of this stuff.

	return nil
}
