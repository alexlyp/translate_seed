// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	_ "net/http/pprof"
	"os"
	"runtime"

	"fmt"

	"github.com/decred/translate_seed/internal/prompt"
)

func main() {
	// Use all processor cores.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Work around defer not working after os.Exit.
	if err := translateMain(); err != nil {
		os.Exit(1)
	}
}

// walletMain is a work-around main function that is required since deferred
// functions (such as log flushing) are not called with calls to os.Exit.
// Instead, main runs this function and checks for a non-nil error, at which
// point any defers have already run, and if the error is non-nil, the program
// can be exited with an error exit status.
func translateMain() error {
	// Load configuration and parse command line.  This function also
	// initializes logging and configures it accordingly.
	reader := bufio.NewReader(os.Stdin)
	seed, err := prompt.Setup(reader)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("Decoded seed hex: %x\n", seed)
	return nil
}
