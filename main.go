package main

import (
	"flag"
	"fmt"
	"github.com/toddkazakov/kyaml-a2c/annotation"
	"log"
)

const ModeAddComments = "add"
const ModeStripComments = "strip"

func main() {
	var path, prefix, mode string
	var strip bool
	flag.StringVar(&path, "path", "", "path to kubernetes manifests")
	flag.StringVar(&prefix, "prefix", annotation.DefaultAnnotationPrefix, "the prefix used in annotations for adding or removing comments")
	flag.StringVar(&mode, "mode", ModeAddComments, "whether to add or strip comments based on the defined annotations. Valid values are 'add' or 'strip'")
	flag.BoolVar(&strip, "strip", false, "strip comments from fields defined in the annotations")
	flag.Parse()

	if path == "" {
		log.Fatalf("path can't be empty")
	}

	var execute func(path, prefix string) error
	switch mode {
	case ModeAddComments:
		execute = annotation.AddCommentsToManifests
	case ModeStripComments:
		execute = annotation.StripCommentsFromManifests
	default:
		log.Fatalf("invalid mode: %v", mode)
	}

	if err := execute(path, prefix); err != nil {
		log.Fatalf("Unexpected error occurred: %v", err.Error())
	}

	fmt.Println("Manifests updated successfully")
}
