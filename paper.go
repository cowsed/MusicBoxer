package main

import "fmt"

type PaperDot struct {
	Pitch, Time float64
}
type OutofBoundsHandleStyle int

const (
	FailOutofBounds OutofBoundsHandleStyle = iota
	SkipOutofBounds
)

//Dots are the music dots with note and time values
//unit ticks is how many midi ticks = 1 unit of time on paper
//vertical unit is how many mm per 1 unit of time on paper
//horizontal unit is how many mm per 1 unit of pitch on paper
func MusicDotToPaperDot(dots []MusicDot, unitTicks, verticalUnit, horizontalUnit float64, OctaveDisplacement int, SkipMode OutofBoundsHandleStyle) ([]PaperDot, error) {
	pdots := make([]PaperDot, len(dots))
	SkippedNotes := 0
	for i, d := range dots {
		note := d.Note
		time := d.Time

		noteName, in := MidiToName[note+uint8(OctaveDisplacement*12)]

		if !in {
			SkippedNotes++
			if SkipMode == FailOutofBounds {
				return nil, fmt.Errorf("error - no note of number %d", note)
			}
		}

		notePos, in := NameToPos[noteName]
		if !in {
			SkippedNotes++
			if SkipMode == FailOutofBounds {
				return nil, fmt.Errorf("error - no note on music box with name %s", noteName)
			}
		}

		dot := PaperDot{
			Pitch: float64(notePos) * horizontalUnit,
			Time:  (float64(time) / unitTicks) * verticalUnit,
		}
		pdots[i] = dot
	}
	fmt.Printf("finished: skipped %d notes\n", SkippedNotes)
	return pdots, nil
}

var NameToPos = map[string]int{
	"C1": 0,
	"D1": 1,
	"G1": 2,
	"A1": 3,
	"B1": 4,

	"C2":  5,
	"D2":  6,
	"E2":  7,
	"F2":  8,
	"F#2": 9,
	"G2":  10,
	"G#2": 11,
	"A2":  12,
	"A#2": 13,
	"B2":  14,

	"C3":  15,
	"C#3": 16,
	"D3":  17,
	"D#3": 18,
	"E3":  19,
	"F3":  20,
	"F#3": 21,
	"G3":  22,
	"G#3": 23,
	"A3":  24,
	"A#3": 25,
	"B3":  26,

	"C4": 27,
	"D4": 28,
	"E4": 29,
}

var MidiToName = map[uint8]string{
	48: "C1",
	49: "C#1",
	50: "D1",
	52: "E1",
	53: "F1",
	54: "F#1",
	55: "G1",
	56: "G#2",
	57: "A1",
	58: "A#1",
	59: "B1",

	48 + 12: "C2",
	49 + 12: "C#2",
	50 + 12: "D2",
	52 + 12: "E2",
	53 + 12: "F2",
	54 + 12: "F#2",
	55 + 12: "G2",
	56 + 12: "G#2",
	57 + 12: "A2",
	58 + 12: "A#2",
	59 + 12: "B2",

	48 + 24: "C3",
	49 + 24: "C#3",
	50 + 24: "D3",
	52 + 24: "E3",
	53 + 24: "F3",
	54 + 24: "F#3",
	55 + 24: "G3",
	56 + 24: "G#3",
	57 + 24: "A3",
	58 + 24: "A#3",
	59 + 24: "B3",

	48 + 36: "C4",
	49 + 36: "C#4",
	50 + 36: "D4",
	52 + 36: "E4",
	53 + 36: "F4",
	54 + 36: "F#4",
	55 + 36: "G4",
	56 + 36: "G#4",
	57 + 36: "A4",
	58 + 36: "A#4",
	59 + 36: "B4",
}
