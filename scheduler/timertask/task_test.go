package timertask

import (
	"fmt"
	"testing"

	"github.com/xfrr/goffmpeg/transcoder"
)

func TestMain(m *testing.M) {
	m.Run()
}

//test goffmpeg work flow
func TestGoffmpeg(t *testing.T) {
	t.Run("begin...", testGo)
}

func testGo(t *testing.T) {
	var inputPath = "/home/pace/go/src/github.com/pili-video-server/videos/test"
	var outputPath = "/home/pace/go/src/github.com/pili-video-server/videos/test360p"

	// Create new instance of transcoder
	trans := new(transcoder.Transcoder)

	// Initialize transcoder passing the input file path and output file path
	err := trans.Initialize(inputPath, outputPath)
	// Handle error...
	if err != nil {
		fmt.Printf("init err :%v!\n", err)
	}

	// SET Resolution TO MEDIAFILE
	// trans.MediaFile().SetResolution("360")

	// SET FRAME RATE TO MEDIAFILE
	trans.MediaFile().SetFrameRate(70)

	// Start transcoder process to check progress
	done := trans.Run(true)

	// Returns a channel to get the transcoding progress
	progress := trans.Output()

	// Example of printing transcoding progress
	for msg := range progress {
		fmt.Println(msg)
	}

	// This channel is used to wait for the transcoding process to end
	err = <-done
}
