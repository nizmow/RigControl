package main

import (
	"flag"
	"log"

	"rigcontrol/internal/machine"
	"rigcontrol/internal/ui"
)

func main() {
	defaultStore, err := machine.DefaultStore()
	if err != nil {
		log.Fatal(err)
	}

	machinesConfig := flag.String("machines-config", defaultStore.Path, "path to machines JSON config")
	flag.Parse()

	store := machine.Store{Path: *machinesConfig}
	profiles, err := store.LoadProfiles()
	if err != nil {
		log.Fatal(err)
	}

	if err := ui.RunWithProfiles(profiles); err != nil {
		log.Fatal(err)
	}
}
