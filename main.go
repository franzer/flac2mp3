package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
)

func mkFolder(folder, bitrate string) string {
	// TODO - determine if linux, windows, macos, etc
	//folType := fmt.Sprintf("%s", bitrate)
	if strings.Contains(folder, "FLAC") {
		chgName := strings.Replace(folder, "FLAC", bitrate, -1)
		newpath := filepath.Join("/", chgName)
		err := os.MkdirAll(newpath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		return newpath
	} else {
		folder_v0 := fmt.Sprintf(`%s %s`, folder, bitrate)
		newpath := filepath.Join("/", folder_v0)
		err := os.MkdirAll(newpath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		return newpath
	}
}

func execLame(wg *sync.WaitGroup, filename, newLoc, bitrate string) {
	// Converts a file from flac to mp3 v0 using FFMPEG //

	defer wg.Done()
	var args []string
	file := filepath.Base(filename)
	mp3out := fmt.Sprintf(`%s/%s.mp3`, newLoc, strings.TrimSuffix(file, ".flac"))
	switch bitrate {
	case "320":
		args = []string{`-y`, `-i`, filename, `-codec:a`, "libmp3lame", `-b:a`, "320k", `-map_metadata`, "0", `-id3v2_version`, "3", `-write_id3v1`, "1", mp3out}
	case "V0":
		args = []string{`-y`, `-i`, filename, `-codec:a`, "libmp3lame", `-q:a`, "0", `-map_metadata`, "0", `-id3v2_version`, "3", `-write_id3v1`, "1", mp3out}
	}
	cmd := exec.Command("ffmpeg", args...)
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
}

func convertFiles(oldFolder, bitrate string) error {
	var wg sync.WaitGroup
	files, err := os.ReadDir(oldFolder)
	if err != nil {
		return err
	}
	newFolder := mkFolder(oldFolder, bitrate)
	newLoc := newFolder

	for _, file := range files {
		path := path.Ext(file.Name())
		if path == ".flac" {
			wg.Add(1)
			newFile := (fmt.Sprintf(`%s/%s`, strings.TrimRight(oldFolder, "/"), file.Name()))
			//fmt.Printf("%s - %s\n", newFile, newLoc)

			go execLame(&wg, newFile, newLoc, bitrate)
		}
	}
	wg.Wait()
	return errors.New(fmt.Sprintf("Files have been saved to:\n%s", newLoc))
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please provide a music folder and bitrate as a command line argument and try again.")
		return
	}

	folder := os.Args[1]
	bitrate := os.Args[2]
	//newFolder := mkFolder(folder)
	s := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
	s.Prefix = (fmt.Sprintf("Converting %s...", filepath.Base(folder)))
	s.Start()
	run := convertFiles(folder, bitrate)
	s.Stop()
	fmt.Println("\n", run)
}
