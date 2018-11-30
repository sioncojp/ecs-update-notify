package main

import (
	"flag"
	"os"

	"github.com/sioncojp/ecs-update-notify"
)

func main() {
	file := flag.String("c", "", "toml file")
	flag.Parse()
	os.Exit(ecsupdatenotify.Run(*file))
}
