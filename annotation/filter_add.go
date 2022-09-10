package annotation

import (
	"sigs.k8s.io/kustomize/kyaml/yaml"
	"strings"
)

// AddCommentsFilter adds annotation values as comment to a field in the yaml document with the path defined in the
// annotation key. The key must be of the form "{AnnotationPrefix}{FieldPath}"
//
// One known limitation is that comments can only be added to scalar fields (e.g. strings, integers, booleans),
// but not to maps and sequences.
//
// Example input:
// ```yaml
// metadata:
//   annotations:
//     comment.kyaml.io/set.spec.values.imageTag: '{"$imagepolicy": "flux-system:podinfo:tag"}'
// spec:
//   values:
//     imageTag: 5.0.0
// ```
//
// Example output (AnnotationPrefix="comment.kyaml.io/set."):
// ```yaml
// metadata:
//   annotations:
//     comment.kyaml.io/set.spec.values.imageTag: '{"$imagepolicy": "flux-system:podinfo:tag"}'
// spec:
//   values:
//     imageTag: 5.0.0 #{"$imagepolicy": "flux-system:podinfo:tag"}
// ```
type AddCommentsFilter struct {
	AnnotationPrefix string
	Callback         func(field, oldValue, newValue string)
}

func (a *AddCommentsFilter) SetCallback(cb func(field, oldValue, newValue string)) {
	a.Callback = cb
}

func (a *AddCommentsFilter) Filter(object *yaml.RNode) (*yaml.RNode, error) {
	return object, a.accept(object)
}

func (a *AddCommentsFilter) accept(object *yaml.RNode) error {
	fieldComments := createFieldCommentMap(a.AnnotationPrefix, object)
	if len(fieldComments) == 0 {
		return nil
	}

	for field, comment := range fieldComments {
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
		node.Document().LineComment = comment
		if a.Callback != nil {
			a.Callback(field, oldValue, comment)
		}
	}

	return nil
}
