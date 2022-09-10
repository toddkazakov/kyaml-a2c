package annotation_test

import (
	"bytes"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.org/toddkazakov/kyaml-a2c/annotation"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	"testing"
)

const withoutComments = `apiVersion: v1
kind: Service
metadata:
  annotations:
    comment.kyaml.io/set.spec.ports: 'ports'
    comment.kyaml.io/set.metadata.name: 'metadata name'
    comment.kyaml.io/set.spec.selector.app: 'selector app name'
  labels:
    app: keto
  name: keto
  namespace: qa2
spec:
  ports:
  - name: http
    port: 8080
    protocol: TCP
    targetPort: 4466
  selector:
    app: keto
  type: ClusterIP
`

const withComments = `apiVersion: v1
kind: Service
metadata:
  annotations:
    comment.kyaml.io/set.spec.ports: 'ports'
    comment.kyaml.io/set.metadata.name: 'metadata name'
    comment.kyaml.io/set.spec.selector.app: 'selector app name'
  labels:
    app: keto
  name: keto # metadata name
  namespace: qa2
spec:
  ports:
  - name: http
    port: 8080
    protocol: TCP
    targetPort: 4466
  selector:
    app: keto # selector app name
  type: ClusterIP
`

func TestAddCommentsFilter(t *testing.T) {
	filter := &annotation.AddCommentsFilter{
		AnnotationPrefix: annotation.DefaultAnnotationPrefix,
	}
	if err := compareFilterResult(filter, withoutComments, withComments); err != nil {
		t.Error(err)
	}
}

func TestStripCommentsFilter(t *testing.T) {
	filter := &annotation.StripCommentsFilter{
		AnnotationPrefix: annotation.DefaultAnnotationPrefix,
	}
	if err := compareFilterResult(filter, withComments, withoutComments); err != nil {
		t.Error(err)
	}
}

func compareFilterResult(filter yaml.Filter, input string, expected string) error {
	parsedInput := yaml.MustParse(input)
	rn, _ := parsedInput.Pipe(filter)

	buf := &bytes.Buffer{}
	err := kio.ByteWriter{Writer: buf}.Write([]*yaml.RNode{rn})
	if err != nil {
		return err
	}

	diff := cmp.Diff(buf.String(), expected)
	if diff != "" {
		return fmt.Errorf("unexpected result: %v", diff)
	}
	return nil
}
