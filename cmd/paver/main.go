package main

import (
	"fmt"
	"os"

	"github.com/farit2000/paver/internal/pkg/app/paver"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		panic(fmt.Sprintf("input shoud be in format ./program <path to dir with packages>, but has %v", args))
	}

	pv := paver.NewPaver(args[1])

	err := pv.Run()
	if err != nil {
		panic(fmt.Sprintf("error while configure workspace: %v", err))
	}
}
