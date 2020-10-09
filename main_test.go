package snd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestWhipSound(t *testing.T) {
	if track, err := LoadTrack("in/2720.dat"); err != nil {
		t.Error(err)
	} else if data, err := track.generate(); err != nil {
		t.Error(err)
	} else if err := ioutil.WriteFile("2720.test", data, 0777); err != nil {
		t.Error(err)
	}
}

func TestSoundsExportIndividual(t *testing.T) {
	var files []string

	err := filepath.Walk("in/", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})

	if err != nil {
		t.Error(err)
	}

	for _, filename := range files {
		func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("uncaught error exporting track %s: %s\n", filename, err)
				}
			}()
			if track, err := LoadTrack(filename); err != nil {
				t.Logf("count not load track %s: %s", filename, err.Error())
			} else {
				_, name := filepath.Split(filename)
				index, err := strconv.Atoi(name[:len(name)-4])

				if err != nil {
					t.Logf("poopoo %s: %s", filename, err.Error())
				}

				raw, err := track.generate()

				if err != nil {
					fmt.Printf("could not generate track %s: %s\n", filename, err.Error())
					return
				}

				if err := ioutil.WriteFile(fmt.Sprintf("out/%d.wav", index), raw, 0777); err != nil {
					fmt.Printf("could not export track %s: %s\n", filename, err.Error())
				}

				t.Log("exported", index)
			}
		}()
	}

}

func TestSoundsExportPacked(t *testing.T) {
	f, err := os.Open("sounds.dat")

	if err != nil {
		t.Error(err)
	}

	tracks, err := LoadTracks(f)

	if err != nil {
		t.Error(err)
	}

	for index, track := range tracks {
		func(index int) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("uncaught error exporting track %d: %s\n", index, err)
				}
			}()
			raw, err := track.generate()

			if err != nil {
				fmt.Printf("could not generate track %d: %s\n", index, err.Error())
				return
			}

			if err := ioutil.WriteFile(fmt.Sprintf("wav/%d.wav", index), raw, 0777); err != nil {
				fmt.Printf("could not save track %d: %s\n", index, err.Error())
			}
		}(index)
	}
}
