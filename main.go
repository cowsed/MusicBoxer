package main

import (
	"fmt"
	"log"
	"os"

	g "github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
)

var dTimePerLine int32 = 240

var (
	BetweenRhythmLinesMM  float32 = 8
	BetweenPitchLinesMM   float32 = 2
	octaveDisplacementAmt int32   = 0
)

var texture *g.Texture
var FirstRun = true

var filename = "No File Selected"
var FileLoaded = false
var CurrentMusicDots []MusicDot
var CurrentMDString string

var TravelHeight float32 = 3
var DrawHeight float32 = -3

var Orientation int32 = 0 //0 is X is time, Y is pitch, 1 is X is ptch, Y is time

var flipX, flipY bool

var XYFeedRate float32 = 2000
var ZFeedRate float32 = 50

func loop() {
	if FirstRun {
		img, _ := g.LoadImage("Paper/blank.png")
		g.EnqueueNewTextureFromRgba(img, func(tex *g.Texture) {
			fmt.Println("Textured")
			texture = tex
		})
	}

	fileControls := g.Layout{
		g.Labelf("File: %s", filename),
		g.Button("Load File").OnClick(loadFileButtonClicked),
	}
	GenerationControls := g.Layout{}
	if FileLoaded {
		GenerationControls = g.Layout{
			g.InputInt(&dTimePerLine).Label("Ticks Per Paper Division"),
			g.Tooltip("The number of midi ticks that will equal 1 unit of space on the paper"),
			g.InputInt(&octaveDisplacementAmt).Label("Octave Displacement"),
			g.Tooltip("To get music to fit in the range of the music box. -1 corresponds to down 1 octave, +1 to up 1 octave"),
			g.Separator(),

			g.InputFloat(&BetweenPitchLinesMM).Label("Pitch Spacing (MM)"),
			g.Tooltip("The Distance between horizontal lines for pitch in millimeters of the music paper. 2 mm on mine"),
			g.InputFloat(&BetweenRhythmLinesMM).Label("Rhythm Spacing (MM)"),
			g.Tooltip("The Distance between vertical lines in millimeters of the music paper. 8 mm on mine"),
			g.Separator(),
			g.InputFloat(&TravelHeight).Label("Travel Height (MM)"),
			g.Tooltip("The height at which the marker will move when travelling between to points"),
			g.InputFloat(&DrawHeight).Label("Draw Height (MM)"),
			g.Tooltip("The height to which the marker will descend to when drawing on the paper"),
			g.Separator(),

			g.Combo("Axes Setup", "X Time", []string{"X Time", "Y Time"}, &Orientation),
			g.Tooltip("Setup of CNC Axes. \nOption 1 - X is time, Y is pitch\nOption 2 - Y is pitch, X is time"),
			g.Checkbox("Flip X", &flipX),
			g.Checkbox("Flip Y", &flipY),
			g.Separator(),
			g.InputFloat(&XYFeedRate).Label("XY Feed Rate (mm/min)"),
			g.Tooltip("THe Feed Rate used when moving between two points"),
			g.InputFloat(&ZFeedRate).Label("Z Feed Rate (mm/min)"),
			g.Tooltip("THe Feed Rate used when moving up/down to draw a dot"),
			g.Separator(),

			g.Button("Create Gcode").OnClick(GenerateGCODEButtonClicked),
			g.Separator(),

			g.Label("Music Points"),
			g.Tooltip("Times and Note values of the piece, ** denote that this note is not playable on the music box"),
			g.InputTextMultiline(&CurrentMDString).Flags(g.InputTextFlagsReadOnly).Size(-1, -1),
		}
	}
	controls := g.Layout{
		fileControls,
		g.Separator(),
		GenerationControls,
	}
	preview := g.Image(texture)

	g.SingleWindow().Layout(

		g.SplitLayout(g.DirectionHorizontal, 300, controls, preview),
	)
	FirstRun = false
}

func main() {
	wnd := g.NewMasterWindow("Canvas", 600, 600, 0)

	wnd.Run(loop)
}

func GenerateGCODEButtonClicked() {
	pdots, err := MusicDotToPaperDot(CurrentMusicDots, float64(dTimePerLine), float64(BetweenRhythmLinesMM), float64(BetweenPitchLinesMM), int(octaveDisplacementAmt), SkipOutofBounds)
	if err != nil {
		panic(err)
	}
	GCodeSrc := PaperDotsToGCODE(pdots, float64(TravelHeight), float64(DrawHeight), 0, 0, Orientation, flipX, flipY, XYFeedRate, ZFeedRate)

	fname, err := dialog.File().Save()
	if err != nil {
		log.Println(err)
		return
	}

	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	f.Write([]byte(GCodeSrc))

	err = f.Close()
	if err != nil {
		panic(err)
	}
}

func loadFileButtonClicked() {
	fname, err := dialog.File().Filter("Midi File", "mid", "midi").Load()
	if err != nil {
		log.Println(err)
		FileLoaded = false
		filename = "invalid path"
		CurrentMusicDots = []MusicDot{}
		CurrentMDString = ""
		return
	}
	filename = fname
	CurrentMusicDots, err = readMidi(fname)
	if err != nil {
		log.Println(err)
		return
	}
	FileLoaded = true
	CurrentMDString = DotsToString(CurrentMusicDots)
}

func DotsToString(mds []MusicDot) string {
	s := ""
	for _, d := range mds {
		noteName, noteOk := MidiToName[d.Note]
		_, boxOk := NameToPos[noteName]
		postfix := ""
		if !boxOk || !noteOk {
			postfix = "**"
		}
		s += fmt.Sprintf("%d  - %s %s\n", d.Time, noteName, postfix)
	}
	return s
}
