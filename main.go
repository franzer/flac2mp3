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

func mkFolder(folder string) string {
	// TODO - determine if linux, windows, macos, etc
	if strings.Contains(folder, "[FLAC]") {
		chgName := strings.Replace(folder, "[FLAC]", "[V0]", -1)
		newpath := filepath.Join("/", chgName)
		err := os.MkdirAll(newpath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		return newpath
	} else {
		folder_v0 := fmt.Sprintf(`%s [V0]`, folder)
		newpath := filepath.Join("/", folder_v0)
		err := os.MkdirAll(newpath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		return newpath
	}
}

func execLame(wg *sync.WaitGroup, filename string, newLoc string) {
	// Converts a file from flac to mp3 v0 using FFMPEG //

	defer wg.Done()
	file := filepath.Base(filename)
	mp3out := fmt.Sprintf(`%s/%s.mp3`, newLoc, strings.TrimSuffix(file, ".flac"))
	args := []string{`-y`, `-i`, filename, `-codec:a`, "libmp3lame", `-q:a`, "0", `-map_metadata`, "0", `-id3v2_version`, "3", `-write_id3v1`, "1", mp3out}
	cmd := exec.Command("ffmpeg", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s", err)
	}
	// fmt.Printf("%s", string(out))
}

func convertFiles(oldFolder string) error {
	var wg sync.WaitGroup
	files, err := os.ReadDir(oldFolder)
	if err != nil {
		return err
	}
	newFolder := mkFolder(oldFolder)
	newLoc := newFolder

	for _, file := range files {
		path := path.Ext(file.Name())
		if path == ".flac" {
			wg.Add(1)
			newFile := (fmt.Sprintf(`%s/%s`, strings.TrimRight(oldFolder, "/"), file.Name()))
			//fmt.Printf("%s - %s\n", newFile, newLoc)

			go execLame(&wg, newFile, newLoc)
		}
	}
	wg.Wait()
	return errors.New(fmt.Sprintf("Files have been saved to:\n%s", newLoc))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a music folder as a command line argument and try again.")
		return
	}

	folder := os.Args[1]
	//newFolder := mkFolder(folder)
	s := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
	s.Prefix = (fmt.Sprintf("Converting %s...", filepath.Base(folder)))
	s.Start()
	run := convertFiles(folder)
	s.Stop()
	fmt.Println("\n", run)
}
