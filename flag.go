package main

import (
	"flag"
)

type Options struct {
	entries int
	wait    int
}

func parse_option() *Options {
	ret := &Options{}
	flag.IntVar(&ret.entries, "d", defaultDumpEntries, "Number of dump entries")
	flag.IntVar(&ret.wait, "w", defaultWaitSecToLogFull, "Seconds to wait")
	flag.Parse()
	return ret
}
