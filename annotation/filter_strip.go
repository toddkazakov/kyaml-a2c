package annotation

import (
	"sigs.k8s.io/kustomize/kyaml/yaml"
	"strings"
)

// StripCommentsFilter removes comments from fields in the yaml document as defined in the annotations.
// This filter is useful for cleaning up comments set by AddCommentsFilter when used in a filter pipe.
//
// Example input:
// ```yaml
// metadata:
//   annotations:
//     comment.kyaml.io/set.spec.values.imageTag: '{"$imagepolicy": "flux-system:podinfo:tag"}'
// spec:
//   values:
//     imageTag: 5.0.0 #{"$imagepolicy": "flux-system:podinfo:tag"}
// ```
//
// Example output (AnnotationPrefix="comment.kyaml.io/set."):
// ```yaml
// metadata:
//   annotations:
//     comment.kyaml.io/set.spec.values.imageTag: '{"$imagepolicy": "flux-system:podinfo:tag"}'
// spec:
//   values:
//     imageTag: 5.0.0
// ```
type StripCommentsFilter struct {
	AnnotationPrefix string
	Callback         func(field, oldValue, newValue string)
}

func (s *StripCommentsFilter) SetCallback(cb func(field, oldValue, newValue string)) {
	s.Callback = cb
}

func (s *StripCommentsFilter) Filter(object *yaml.RNode) (*yaml.RNode, error) {
	return object, s.accept(object)
}

func (s *StripCommentsFilter) accept(object *yaml.RNode) error {
	fieldComments := createFieldCommentMap(s.AnnotationPrefix, object)
	if len(fieldComments) == 0 {
		return nil
	}

	for field, _ := range fieldComments {
		// Get the field node given a field path (e.g. "spec.template.spec.container.[name=nginx]")
		path := strings.Split(field, ".")
		node, err := object.Pipe(yaml.Lookup(path...))

		// Return error if a lookup fails
		if err != nil {
			return err
		}

		// TODO: figure out how to set comments to Mapping Or Sequence nodes correctly
		if node.IsNilOrEmpty() || node.YNode().Kind != yaml.ScalarNode {
			continue
		}

		oldValue := node.Document().LineComment
		node.Document().LineComment = ""
		if s.Callback != nil {
			s.Callback(field, oldValue, "")
		}
	}

	return nil
}
