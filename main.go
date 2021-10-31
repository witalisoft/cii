package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const version = "0.1.0"

type config struct {
	image    *string
	nocolor  *bool
	noformat *bool
	platform *string
	version  *bool
}

func main() {
	config := getCliFlags()
	if *config.version {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}
	if *config.image == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	imageManifest, err := config.getImageManifest()
	if err != nil {
		log.Fatalf("problem with get image manifest, err: %v", err)
	}
	info, err := config.writeImageInfo(&imageManifest)
	if err != nil {
		log.Fatalf("problem with writing image info, err: %v", err)
	}
	info.WriteTo(os.Stdout)
	imageHistory, err := config.getImageHistory()
	if err != nil {
		log.Fatalf("cannot get image history, err: %v", err)
	}
	err = config.prepareImageHistory(&imageHistory)
	if err != nil {
		log.Fatalf("problem with preparing image history, err: %v", err)
	}
	parsedPlatform, err := parsePlatform(config.platform)
	if err != nil {
		log.Fatalf("problem with parsing plaform parameter, err: %v", err)
	}
	imagePlatformManifest, err := config.getImagePlaformManifest(parsedPlatform)
	if err != nil {
		log.Fatalf("problem with get image platform manifest, err: %v", err)
	}
	config.updateLayerSize(&imageHistory, &imagePlatformManifest)
	history, err := config.writeHistory(&imageHistory)
	if err != nil {
		log.Fatalf("problem with writing image history, err: %v", err)
	}
	history.WriteTo(os.Stdout)
}

func getCliFlags() *config {
	config := config{
		image: flag.String("image", "", "image name"),
		platform: flag.String("platform", "linux/amd64",
			"specify platform of which you want to get layers history in the form os/arch"),
		nocolor:  flag.Bool("no-color", false, "disable color output"),
		noformat: flag.Bool("no-format", false, "don't try to format shell scripts in the image history"),
		version:  flag.Bool("version", false, "print the version and exit"),
	}
	flag.Parse()
	return &config
}
