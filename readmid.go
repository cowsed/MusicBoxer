package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/algoGuy/EasyMIDI/smf"
	"github.com/algoGuy/EasyMIDI/smfio"
)

type MusicDot struct {
	Time int
	Note uint8
}

func readMidi(fname string) ([]MusicDot, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read and save midi to smf.MIDIFile struct
	midi, err := smfio.Read(bufio.NewReader(file))

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error parsing file")
	}

	var Notes = []MusicDot{}
	for i := 0; i < int(midi.GetTracksNum()); i++ {

		// Get zero track
		track := midi.GetTrack(uint16(i))

		// Get all midi events via iterator
		iter := track.GetIterator()

		var TotalTime int = 0
		for iter.MoveNext() {
			val := iter.GetValue()
			//Keep track of time
			TotalTime += int(val.GetDTime())

			//Chck if it's not one
			status := val.GetStatus()
			if status != smf.NoteOnStatus {
				//Dont need it, continue
				continue
			}

			data := val.GetData()

			note := data[0]
			vel := data[1]
			if vel == 0 {
				//0 sound
				continue
			}
			Notes = append(Notes, MusicDot{Time: TotalTime, Note: note})

		}
	}
	return Notes, nil
}
