package main

import (
	"fmt"
	"strings"
)

var DefualtDownPositionMM = -1
var DefualtUpPositionMM = 2

func PaperDotsToGCODE(dots []PaperDot, UpPositionMM, DownPositionMM float64, XOffset, YOffset float64, orientation int32, Xflip, Yflip bool, XYFeed, ZFeed float32) string {
	b := strings.Builder{}

	downLine := fmt.Sprintf("G0 Z%f F%f\n", DownPositionMM, ZFeed)
	upLine := fmt.Sprintf("G0 Z%f F%f\n", UpPositionMM, ZFeed)

	//Write GCODE Header
	b.WriteString("G21\n")
	b.WriteString(upLine)

	//p.X is pitch, p.Y is time
	for _, p := range dots {
		var machineX, machineY float64
		//0 is X is time, Y is pitch, 1 is X is ptch, Y is time
		if orientation == 0 {
			machineX = p.Time  // Time
			machineY = p.Pitch // Pitch
		} else {
			machineX = p.Pitch //Pitch
			machineY = p.Time  //Time

		}
		if flipX {
			machineX *= -1
		}
		if flipY {
			machineY *= -1
		}

		gotoLine := fmt.Sprintf("G0 X%f Y%f F%.f\n", machineX, machineY, XYFeed)

		b.WriteString(gotoLine)
		b.WriteString(downLine)
		b.WriteString(upLine)
	}
	return b.String()
}
