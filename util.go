package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	units "github.com/docker/go-units"
	"github.com/fatih/color"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

type ImagePlatformManifest struct {
	Layers []struct {
		Size int `json:"size"`
	} `json:"layers"`
}

type ImageManifest struct {
	Manifests []struct {
		Digest   string        `json:"digest"`
		Platform ImagePlatform `json:"platform"`
	} `json:"manifests"`
}

type ImagePlatform struct {
	Architecture string `json:"architecture"`
	Os           string `json:"os"`
	Variant      string `json:"variant,omitempty"`
}

func (c *config) getImagePlaformManifest(platform *v1.Platform) (ImagePlatformManifest, error) {
	manifest, err := crane.Manifest(*c.image, crane.WithPlatform(platform))
	if err != nil {
		return ImagePlatformManifest{}, err
	}
	imagePlatformManifest := ImagePlatformManifest{}
	err = json.Unmarshal(manifest, &imagePlatformManifest)
	if err != nil {
		return ImagePlatformManifest{}, err
	}
	return imagePlatformManifest, nil
}

func (c *config) getImageManifest() (ImageManifest, error) {
	manifest, err := crane.Manifest(*c.image)
	if err != nil {
		return ImageManifest{}, err
	}
	imageManifest := ImageManifest{}
	err = json.Unmarshal(manifest, &imageManifest)
	if err != nil {
		return ImageManifest{}, err
	}
	return imageManifest, nil
}

func (c *config) writeImageInfo(imageManifest *ImageManifest) (bytes.Buffer, error) {
	var out bytes.Buffer
	color := color.New(color.FgHiWhite, color.Bold)
	if *c.nocolor {
		color.DisableColor()
	}
	color.Fprint(&out, "Available platforms:\n\n")
	color.Fprintf(&out, "%17s\t%8s\t%71s\n", "platform", "size", "digest")
	for _, value := range imageManifest.Manifests {
		platformString := platformToString(value.Platform)
		parsedPlatform, err := parsePlatform(&platformString)
		if err != nil {
			return out, err
		}
		imagePlatformManifest, err := c.getImagePlaformManifest(parsedPlatform)
		if err != nil {
			return out, err
		}
		fmt.Fprintf(&out, "%17s\t%8s\t%71s\n",
			platformString,
			imagePlatformSize(&imagePlatformManifest),
			value.Digest)
	}
	return out, nil
}

func parsePlatform(platform *string) (*v1.Platform, error) {
	p := &v1.Platform{}
	parts := strings.Split(*platform, "/")

	if len(parts) < 2 {
		return nil, fmt.Errorf("failed to parse platform '%s': expected format os/arch[/variant]", *platform)
	}
	if len(parts) > 3 {
		return nil, fmt.Errorf("failed to parse platform '%s': too many slashes", *platform)
	}

	p.OS = parts[0]
	p.Architecture = parts[1]
	if len(parts) > 2 {
		p.Variant = parts[2]
	}

	return p, nil
}

func platformToString(platform ImagePlatform) string {
	p := ""
	if platform.Os != "" && platform.Architecture != "" {
		p = platform.Os + "/" + platform.Architecture
	}
	if platform.Variant != "" {
		p += "/" + platform.Variant
	}
	return p
}

func imagePlatformSize(imagePlatformManifest *ImagePlatformManifest) string {
	sum := 0
	for _, value := range imagePlatformManifest.Layers {
		sum += value.Size
	}
	return units.BytesSize(float64(sum))
}
