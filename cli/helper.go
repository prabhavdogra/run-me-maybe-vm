package cli

import (
	"fmt"
	"os"
)

func GetArgs() Args {
	args := Args{}
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <sourcefile.cmm>\n", os.Args[0])
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
