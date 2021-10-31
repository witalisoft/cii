package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	units "github.com/docker/go-units"
	"github.com/fatih/color"
	"github.com/google/go-containerregistry/pkg/crane"
	"mvdan.cc/sh/v3/syntax"
)

type ImageHistory struct {
	History []struct {
		Created          string `json:"created"`
		Created_by       string `json:"created_by"`
		Empty_layer      bool   `json:"empty_layer,omitempty"`
		Layer_size       string
		ShellFormatError error
	} `json:"history"`
	Created string `json:"created"`
}

func (c *config) getImageHistory() (ImageHistory, error) {
	cfg, err := crane.Config(*c.image)
	if err != nil {
		return ImageHistory{}, err
	}
	imageHistory := ImageHistory{}
	err = json.Unmarshal(cfg, &imageHistory)
	if err != nil {
		return ImageHistory{}, err
	}
	return imageHistory, nil
}

func (c *config) writeHistory(imageHistory *ImageHistory) (bytes.Buffer, error) {
	var out bytes.Buffer
	color := color.New(color.FgHiWhite, color.Bold)
	if *c.nocolor {
		color.DisableColor()
	}
	dataLayerCnt := 0
	emptyLayerCnt := 0
	for _, value := range imageHistory.History {
		if value.Empty_layer {
			emptyLayerCnt++
		} else {
			dataLayerCnt++
		}
	}
	color.Fprint(&out, "\n\nData layers: ")
	fmt.Fprintf(&out, "%d\n", dataLayerCnt)
	color.Fprint(&out, "Empty layers: ")
	fmt.Fprintf(&out, "%d\n", emptyLayerCnt)
	createdTime, err := time.Parse(time.RFC3339, imageHistory.Created)
	if err != nil {
		return out, err
	}
	color.Fprint(&out, "Last pushed: ")
	fmt.Fprintf(&out, "%s\n", units.HumanDuration(time.Since(createdTime)))
	color.Fprint(&out, "\n\nLayers history for platform: ")
	fmt.Fprintf(&out, "%s\n\n", *c.platform)
	for i, value := range imageHistory.History {
		color.Fprint(&out, "layer: ")
		fmt.Fprintf(&out, "%d\n", i+1)
		color.Fprint(&out, "size: ")
		fmt.Fprintf(&out, "%s\n", value.Layer_size)
		color.Fprint(&out, "empty_layer: ")
		fmt.Fprintf(&out, "%t\n", value.Empty_layer)
		color.Fprint(&out, "created: ")
		fmt.Fprintf(&out, "%s\n", value.Created)
		if imageHistory.History[i].ShellFormatError != nil {
			color.Fprint(&out, "shell format error: ")
			fmt.Fprintf(&out, "%s\n", imageHistory.History[i].ShellFormatError.Error())
		}
		color.Fprint(&out, "created_by: ")
		if strings.Contains(value.Created_by, "\n") {
			fmt.Fprint(&out, "|\n")
			for _, line := range strings.Split(value.Created_by, "\n") {
				fmt.Fprintf(&out, "\t%s\n", line)
			}
			fmt.Fprint(&out, "\n")
		} else {
			fmt.Fprintf(&out, "%s\n\n", value.Created_by)
		}
	}
	return out, nil
}

func (c *config) prepareImageHistory(imageHistory *ImageHistory) error {
	for i, value := range imageHistory.History {
		createdTime, err := time.Parse(time.RFC3339, value.Created)
		if err != nil {
			return err
		}
		imageHistory.History[i].Created = units.HumanDuration(time.Since(createdTime))
		if !strings.Contains(value.Created_by, "#(nop)") {
			if !*c.noformat {
				shellFmt, err := shellFormatter(value.Created_by)
				if err != nil {
					imageHistory.History[i].ShellFormatError = err
				}
				imageHistory.History[i].Created_by = shellFmt
			}
		}
	}
	return nil
}

func (c *config) updateLayerSize(imageHistory *ImageHistory, imagePlatformManifest *ImagePlatformManifest) {
	dataLayerCnt := 0
	for i, value := range imageHistory.History {
		if !value.Empty_layer {
			imageHistory.History[i].Layer_size =
				units.BytesSize(float64(imagePlatformManifest.Layers[dataLayerCnt].Size))
			dataLayerCnt++
		} else {
			imageHistory.History[i].Layer_size = "0B"
		}
	}
}

func shellFormatter(oper string) (string, error) {
	shellPath := regexp.MustCompile(`\s+`).Split(oper, 2)
	shellExec := filepath.Base(shellPath[0])
	allowedShells := regexp.MustCompile(`sh|bash`)
	if allowedShells.Match([]byte(shellExec)) {
		var parsedOutput bytes.Buffer
		reader := strings.NewReader(oper)
		parser, err := syntax.NewParser().Parse(reader, "")
		if err != nil {
			return oper, err
		}
		syntax.NewPrinter().Print(&parsedOutput, parser)
		return strings.TrimSuffix(parsedOutput.String(), "\n"), nil
	}
	return oper, nil
}
