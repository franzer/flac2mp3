package main

import (
	"fmt"
	"path/filepath"
	"testing"
)

func BenchmarkConversion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// The code to be benchmarked
		folder := `<src_folder>`
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
