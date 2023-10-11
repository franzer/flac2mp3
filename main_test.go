package main

import (
	"fmt"
	"path/filepath"
	"testing"
)

func BenchmarkConversion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// The code to be benchmarked
		//folder := "/home/franz/Downloads/Bakemonogatari Portable Special Content CD [FLAC]"
		folder := `/home/franz/Downloads/[1990.04.05] V.A. - Perfect Collection Dragon Slayer The Legend of Heroes {KICA-1003~4}`
		bitrate := "V0"
		curFileList := []string{}
		newFileList := []string{}
		newFolder := mkFolder(folder, bitrate)

		err := filepath.Walk(folder, visit(folder, newFolder, bitrate, &curFileList, &newFileList))
		if err != nil {
			fmt.Println(err)
		}
	}
}
