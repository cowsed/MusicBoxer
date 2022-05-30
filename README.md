#  Music Boxer
A program to convert a Midi file to the GCODE for a pen plotter

## How To

### Generating the Program
1. Load a Midi file
2. Examine the *Music Points* area for the desired subdivision
3. Select the octave displacement - shouldn't be needed if the notes are setup when composing/arranging
4. Select the Pitch and Rhythm spacing. These are measured from the paper you wish to draw on
5. Adjust the travel and draw height. Travel Height is the height the pen will be at when moving between steps. Draw Height is the height to which the head will descend to draw. These depend on your machine
6. Select which axis of the plotter/Cnc machine will correspond to time on the music paper.
7. Flip X or Flip Y as needed for your machine
8. Adjust the XY and Z feedrate to what is comfortable for your machine
9. Press the generate gcode button and save the file wherever you choose

### Running the Program
1. Square music box paper to the machine
2. Zero the machine at the lowest note at the top of the paper.