package main

import (
	"fmt"
	"os"

	"github.com/kinghanzala/gcpsec/cmd"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	build := cmd.BuildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	}

	if err := cmd.NewRootCmd(build).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
