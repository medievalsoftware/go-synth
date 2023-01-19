package main

import (
	"fmt"
	snd "github.com/thedaneeffect/rs2-snd"
	"os"
)

func main() {
	tracks, err := snd.LoadTracks("sounds.dat")

	if err != nil {
		panic(err)
	}

	if err = os.Mkdir("sounds", 0666); err != nil {
		panic(err)
	}

	for _, track := range tracks {
		var data []byte
		if data, err = track.CreateRiff(); err != nil {
			panic(fmt.Errorf("track riff %d: %v", track.ID, err))
		} else if err = os.WriteFile(fmt.Sprintf("sounds/%d.wav", track.ID), data, 0666); err != nil {
			panic(fmt.Errorf("track save %d: %v", track.ID, err))
		}
	}
}
