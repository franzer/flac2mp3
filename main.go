package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
)

func mkFolder(folder, bitrate string) string {
	// TODO - determine if linux, windows, macos, etc

	folder = filepath.Clean(folder)

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
	// Converts a file from source format to mp3 using FFMPEG //

	var args []string
	defer wg.Done()
	switch bitrate {
	case "320":
		args = []string{`-y`, `-i`, filename, `-codec:a`, "libmp3lame", `-b:a`, "320k", `-map_metadata`, "0", `-id3v2_version`, "3", `-write_id3v1`, "1", newLoc}
	case "V0":
		args = []string{`-y`, `-i`, filename, `-codec:a`, "libmp3lame", `-q:a`, "0", `-map_metadata`, "0", `-id3v2_version`, "3", `-write_id3v1`, "1", newLoc}
	}
	cmd := exec.Command("ffmpeg", args...)
	_, err := cmd.Output()
	if err != nil {
		fmt.Print(".")
	}
}

func visit(oldFolder, newFolder string, bitrate string, curFileList, newFileList *[]string) filepath.WalkFunc {
	// Takes in a target location, walks through the folders and files, and puts them in to their respective slices //
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		relativePath, err := filepath.Rel(oldFolder, path)
		if err != nil {
			return err
		}
		if relativePath == "." {
			return nil
		}
		newPath := filepath.Join(newFolder, relativePath)
		if info.IsDir() {
			return os.MkdirAll(newPath, info.Mode())
		}
		*curFileList = append(*curFileList, path)
		*newFileList = append(*newFileList, strings.Replace(newPath, filepath.Ext(newPath), ".mp3", -1))
		return nil
	}
}

func main() {
	var curFileList []string
	var newFileList []string
	var wg sync.WaitGroup

	folder := flag.String("f", "", "Folder you would like to convert (/home/user/Music/Music Folder)")
	bitrate := flag.String("b", "V0", "Bitrate you want to convert to 320 (CBR) or V0(VRR).")
	flag.Parse()

	if len(os.Args) < 3 {
		fmt.Println("Please provide a music folder and bitrate as a command line argument and try again.  Use -h or --help to see flags.")
		return
	}

	//folder := os.Args[1]
	//bitrate := os.Args[2]
	newFolder := mkFolder(*folder, *bitrate)

	s := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
	s.Prefix = (fmt.Sprintf("Converting %s...", filepath.Base(*folder)))
	s.Start()
	err := filepath.Walk(*folder, visit(*folder, newFolder, *bitrate, &curFileList, &newFileList))
	if err != nil {
		fmt.Println(err)
	}
	for index, file := range curFileList {
		wg.Add(1)
		go execLame(&wg, file, newFileList[index], *bitrate)
	}
	wg.Wait()
	fmt.Printf("Your files are located at:\n%s", newFolder)
	s.Stop()
}
