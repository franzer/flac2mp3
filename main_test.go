package main

import (
	"fmt"
	"path/filepath"
	"testing"
)

func BenchmarkConversion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// The code to be benchmarked
		folder := "/Users/franzer/Downloads/Brave Story Original Soundtrack[FLAC]"
		bitrate := "V0"
		newFolder := mkFolder(folder, bitrate)

		err := filepath.Walk(folder, visit(folder, newFolder, bitrate))
		if err != nil {
			fmt.Println(err)
		}
	}
}
