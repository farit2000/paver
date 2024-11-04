package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/farit2000/paver/internal/pkg/app/paver"
)

func main() {
	workersNum := flag.Int("workers", 2, "an int")
	workDir := flag.String("workdir", "./", "a string")
	if *workersNum <= 0 {
		panic(fmt.Sprintf("workers shoud be more that zero, but: %d", *workersNum))
	}
	flag.Parse()

	pv := paver.NewPaver(*workDir)

	now := time.Now()
	defer func() {
		fmt.Println("Elapsed time:", time.Since(now))
	}()
	err := pv.Run(*workersNum)
	if err != nil {
		panic(fmt.Sprintf("error while configure workspace: %v", err))
	}
}
