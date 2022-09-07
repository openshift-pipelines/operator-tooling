package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	config string
)

func main() {
	flag.StringVar(&config, "config", "", "components configuration to load")

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		os.Exit(1)
	}
	var err error
	switch args[0] {
	case "component-version":
		err = componentVersion(config, args[1:])
	case "check":
		err = check(config, false)
	case "check-bugfix":
		err = check(config, true)
	case "bump":
		err = bump(config, false)
	case "bump-bugfix":
		err = bump(config, true)
	case "generate-upstream-sources":
		err = generateUpstreamSources(config)
	default:
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
