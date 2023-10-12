**flac2MP3**

A FFMPEG wrapper written in Go to quickly convert FLAC files to MP3 in either 320k (CBR) or V0 (VBR) to slim file size and store elsewhere.  This tool utilizes go's concurrency to quickly convert files without needing to understand all of the FFMPEG flags that are used.  Simply add a source folder and the bitrate desired and let it run.  This will maintain all metadata on the files as well.

**How to use:**

```
- Download the latest tool or clone the repo and build with go build .
- run flac2mp3 -f <folder_location> -b <bitrate>
  - ex. flac2mp3 -f /home/User/Music/Folder -b V0
- Wait for flac2mp3 to finish
- Your new folder will be in /home/User/Music/Folder V0
```
