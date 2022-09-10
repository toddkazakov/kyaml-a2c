package annotation

import (
	"fmt"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

const DefaultAnnotationPrefix = "comment.kyaml.io/set."

func AddCommentsToManifests(packagePath, annotationPrefix string) error {
	rw := &kio.LocalPackageReadWriter{
		KeepReaderAnnotations: false,
		PreserveSeqIndent:     true,
		PackagePath:           packagePath,
		NoDeleteFiles:         true,
	}
	pipeline := kio.Pipeline{
		Inputs:  []kio.Reader{rw},
		Outputs: []kio.Writer{rw},
		Filters: []kio.Filter{
			NewAddCommentsFilter(annotationPrefix, addCommentsCallback),
		},
	}
	return pipeline.Execute()
}

func StripCommentsFromManifests(packagePath, annotationPrefix string) error {
	rw := &kio.LocalPackageReadWriter{
		KeepReaderAnnotations: false,
		PreserveSeqIndent:     true,
		PackagePath:           packagePath,
		NoDeleteFiles:         true,
	}
	pipeline := kio.Pipeline{
		Inputs:  []kio.Reader{rw},
		Outputs: []kio.Writer{rw},
		Filters: []kio.Filter{
			NewStripCommentsFilter(annotationPrefix, stripCommentsCallback),
		},
	}
	return pipeline.Execute()
}

func addCommentsCallback(file, field, comment string, node *yaml.RNode) {
	fmt.Printf("%v: setting %v to %v\n", file, field, comment)
}

func stripCommentsCallback(file, field, comment string, node *yaml.RNode) {
	fmt.Printf("%v: removing comment from %v\n", file, field)
}
