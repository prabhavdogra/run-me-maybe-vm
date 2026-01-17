package main

import (
	"fmt"
	"os"
)

type Args struct {
	FileName  string
	DebugMode bool
}

func parseArgs() Args {
	args := Args{}
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <sourcefile.wm>\n", os.Args[0])
		return args
	}

	args.FileName = os.Args[1]
	for _, arg := range os.Args[2:] {
		if arg == "--debug" || arg == "-d" {
			args.DebugMode = true
		}
	}
	return args
}
