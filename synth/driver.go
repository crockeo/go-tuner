package synth

import (
	"fmt"
	"github.com/HardWareGuy/portaudio-go"
	"github.com/crockeo/go-tuner/config"
	"math"
	"time"
)

// Type Driver is an interface to define the required behavior for a data
// structure that can be used in driving PortAudio music synthesis.
type Driver interface {
	// Calculating the time in seconds that this driver should be running.
	CalculateDuration() time.Duration

	// Getting the number of output channels this driver is expecting.
	OutputChannels() int

	// Calculating the output on whatever set of channels for a given driver.
	CalculateOutput() []float32

	// Stepping the internal phases given a sample rate.
	StepPhases(int)
}

// The driver to play a single note for a given duration.
type SingleDriver struct {
	Note      NoteData // The data for this note.
	Phase     float32  // The current phase of the driver.
	Time      float32  // The current time of the SingleDriver.
	StartTime float32  // The time this driver was started - if used without a PrimaryDriver it will always be 0.
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
func NewSingleDriverChild(nd NoteData, startTime float32) *SingleDriver {
	sd := new(SingleDriver)

	sd.Note = nd
	sd.Phase = 0
	sd.Time = startTime
	sd.StartTime = startTime

	return sd
}

// Calculating the time in seconds that this driver should be running.
func (sd *SingleDriver) CalculateDuration() time.Duration {
	return time.Duration(sd.Note.Duration)
}

// Getting the number of output channels this driver is expecting.
func (sd *SingleDriver) OutputChannels() int {
	return 2
}

// Calculating the output on whatever set of channels for a given driver.
func (sd *SingleDriver) CalculateOutput() []float32 {
	output := float32(math.Sin(float64(sd.Phase))) * sd.Note.Volume * sd.Note.FadeFunc(sd.Time-sd.StartTime, sd.Note.Duration)

	outputs := make([]float32, sd.OutputChannels())
	for k, _ := range outputs {
		outputs[k] = output
	}

	return outputs
}

// Stepping the internal phases given a sample rate.
func (sd *SingleDriver) StepPhases(sampleRate int) {
	sd.Phase += sd.Note.Frequency * (1 / float32(sampleRate))
	if sd.Phase >= 2*math.Pi {
		sd.Phase -= 2 * math.Pi
	}

	sd.Time += 1.0 / float32(sampleRate)
}

// The primary driver that is used by the rest of the program by default to
// start whichever synth.
type PrimaryDriver struct {
	QueuedNotes  []DelayedNoteData // The list of NoteDatas to add.
	CurrentNotes []*SingleDriver   // The list of current SingleDrivers.
	Time         float32           // The current time of the PrimaryDriver.
	LastTime     float32           // The time that the last SingleDriver was added.
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

// Calculating the time in seconds that this driver should be running.
func (pd *PrimaryDriver) CalculateDuration() time.Duration {
	var delay, max float32

	delay = 0
	max = 0

	for _, qn := range pd.QueuedNotes {
		delay += qn.Delay
		v := delay + qn.ND.Duration

		if v > max {
			max = v
		}
	}

	return time.Duration(max)
}

// Trying to add a new DelayedNoteData to the list of queued notes inside of a
// PrimaryDriver.
func (pd *PrimaryDriver) AddDelayedNote(dnd DelayedNoteData) {
	if config.DebugMode {
		fmt.Print("Adding note: ")
		fmt.Println(dnd.ND)
	}

	if len(pd.QueuedNotes) == 0 {
		pd.LastTime = pd.Time
	}

	pd.QueuedNotes = append(pd.QueuedNotes, dnd)
}

// Getting the number of output channels this driver is expecting.
func (pd *PrimaryDriver) OutputChannels() int { return 2 }

// Calculating the output on whatever set of channels for a given driver.
func (pd *PrimaryDriver) CalculateOutput() []float32 {
	vs := []float32{0.0, 0.0}
	for _, sd := range pd.CurrentNotes {
		for i, o := range sd.CalculateOutput() {
			vs[i] += o
		}
	}

	return vs
}

// Stepping the internal phases given a sample rate.
func (pd *PrimaryDriver) StepPhases(sampleRate int) {
	// Appending new notes to the set of current notes.
	for len(pd.QueuedNotes) > 0 && pd.Time-pd.LastTime <= pd.QueuedNotes[0].Delay {
		if config.DebugMode {
			fmt.Print("Playing note: ")
			fmt.Println(pd.QueuedNotes[0].ND)
		}

		pd.CurrentNotes = append(pd.CurrentNotes, NewSingleDriverChild(pd.QueuedNotes[0].ND, pd.Time))
		pd.QueuedNotes = pd.QueuedNotes[1:]

		pd.LastTime = pd.Time
	}

	// TODO: Remove elements from CurrentNotes that are terminated.

	// Stepping the phases for the sub drivers.
	for _, sd := range pd.CurrentNotes {
		sd.StepPhases(sampleRate)
	}

	pd.Time += 1.0 / float32(sampleRate)
}

// Returning a function to drive music synthesis given a driver and a sample
// rate.
func DriverFunction(driver Driver, sampleRate int) func([][]float32) {
	return func(out [][]float32) {
		for i := range out[0] {
			output := driver.CalculateOutput()
			for j := range output {
				out[j][i] = output[j]
			}

			driver.StepPhases(sampleRate)
		}
	}
}

// Running a synth server.
func RunSynth(driver Driver, errChannel chan error, exitChannel chan bool) {
	err := portaudio.Initialize()
	if err != nil {
		errChannel <- err
		<-exitChannel
	}
	defer portaudio.Terminate()

	stream, err := portaudio.OpenDefaultStream(0, 2, 44100, 0, DriverFunction(driver, 44100))
	if err != nil {
		errChannel <- err
		<-exitChannel
	}
	defer stream.Close()

	stream.Start()
	defer stream.Stop()

	<-exitChannel
}
