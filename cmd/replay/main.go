package main

import "os"

func main() {

	replay := newReplay()

	err := replay()
	if err != nil {
		os.Exit(1)
	}
}

type replayer func() error

func newReplay() replayer {
	return func() error {
		return nil
	}
}
