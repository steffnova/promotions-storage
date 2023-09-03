package csv

import "fmt"

var (
	ErrStreamer  = fmt.Errorf("error setting up csv streamer")
	ErrStreaming = fmt.Errorf("error processing csv file")
)
