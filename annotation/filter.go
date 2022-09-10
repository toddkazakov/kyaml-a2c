package annotation

import (
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/sets"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	"strings"
)

func NewAddCommentsFilter(annotationPrefix string, callback func(file, field, comment string, node *yaml.RNode)) kio.Filter {
	filter := &AddCommentsFilter{
		AnnotationPrefix: annotationPrefix,
	}

	return HandleComments(filter, callback)
}

func NewStripCommentsFilter(annotationPrefix string, callback func(file, field, comment string, node *yaml.RNode)) kio.Filter {
	filter := &StripCommentsFilter{
		AnnotationPrefix: annotationPrefix,
	}

	return HandleComments(filter, callback)
}

func HandleComments(filter FilterWithCallback, callback func(file, field, comment string, node *yaml.RNode)) kio.Filter {
	return kio.FilterFunc(
		func(nodes []*yaml.RNode) ([]*yaml.RNode, error) {
			filesToUpdate := sets.String{}
			for i := range nodes {
				path, _, err := kioutil.GetFileAnnotations(nodes[i])
				if err != nil {
					return nil, err
				}

				filter.SetCallback(func(field, oldValue, newValue string) {
					if oldValue != newValue {
						callback(path, field, newValue, nodes[i])
						filesToUpdate.Insert(path)
					}
				})
				_, err = filter.Filter(nodes[i])
				if err != nil {
					return nil, err
				}
			}

			var nodesInUpdatedFiles []*yaml.RNode
			for i := range nodes {
				path, _, err := kioutil.GetFileAnnotations(nodes[i])
				if err != nil {
					return nil, err
				}
				if filesToUpdate.Has(path) {
					nodesInUpdatedFiles = append(nodesInUpdatedFiles, nodes[i])
				}
			}
			return nodesInUpdatedFiles, nil
		})
}

type FilterWithCallback interface {
	Filter(object *yaml.RNode) (*yaml.RNode, error)
	SetCallback(cb func(field, oldValue, newValue string))
}

func createFieldCommentMap(annotationPrefix string, node *yaml.RNode) map[string]string {
	// Early return if the filter is not initialized correctly
	if len(annotationPrefix) == 0 {
		return nil
	}

	annotations := make(map[string]string, 0)
	for k, v := range node.GetAnnotations() {
		if strings.HasPrefix(k, annotationPrefix) {
			field := strings.TrimPrefix(k, annotationPrefix)
			if field != "" {
				annotations[field] = v
			}
		}
	}

	return annotations
}
