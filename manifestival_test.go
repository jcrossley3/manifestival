package manifestival_test

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"

	. "github.com/kabanero-io/manifestival"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestFromReader(t *testing.T) {
	tests := []struct {
		name                string
		reader              io.Reader
		expectedApiVersions []string
	}{{
		name: "from_bytes",
		reader: bytes.NewReader([]byte(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3
---
apiVersion: v1
kind: Service
spec:
  selector:
    app: MyApp
  ports:
    - protocol: TCP
      port: 80
      targetPort: 9376
`)),
		expectedApiVersions: []string{"apps/v1", "v1"},
	}}

	cl, _ := client.New(nil, client.Options{})
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s", tc.name), func(t *testing.T) {
			m, err := FromReader(tc.reader, cl)
			if err != nil {
				t.Fatalf("FromReader returned: %v", err)
			}

			foundApiVersions := make([]string, 0)
			for _, r := range m.Resources {
				foundApiVersions = append(foundApiVersions, r.GetAPIVersion())
			}
			if !reflect.DeepEqual(tc.expectedApiVersions, foundApiVersions) {
				t.Fatalf("Expected API kinds %v but found %v", tc.expectedApiVersions, foundApiVersions)
			}
		})
	}
}

func TestUpdateChanges(t *testing.T) {
	tests := []struct {
		name    string
		changed bool
		src     map[string]interface{}
		tgt     map[string]interface{}
		want    map[string]interface{}
	}{{
		name:    "identical maps",
		changed: false,
		src: map[string]interface{}{
			"x": map[string]interface{}{
				"y": 1,
			},
		},
		tgt: map[string]interface{}{
			"x": map[string]interface{}{
				"y": 1,
			},
		},
		want: map[string]interface{}{
			"x": map[string]interface{}{
				"y": 1,
			},
		},
	}, {
		name:    "add nested map entry",
		changed: true,
		src: map[string]interface{}{
			"x": map[string]interface{}{
				"y": 1,
			},
		},
		tgt: map[string]interface{}{
			"x": map[string]interface{}{
				"a": "foo",
			},
		},
		want: map[string]interface{}{
			"x": map[string]interface{}{
				"y": 1,
				"a": "foo",
			},
		},
	}, {
		name:    "change nested map entry",
		changed: true,
		src: map[string]interface{}{
			"x": map[string]interface{}{
				"y": 1,
			},
		},
		tgt: map[string]interface{}{
			"x": map[string]interface{}{
				"y": 2,
			},
		},
		want: map[string]interface{}{
			"x": map[string]interface{}{
				"y": 1,
			},
		},
	}, {
		name:    "change missing map entry",
		changed: true,
		src: map[string]interface{}{
			"x": map[string]interface{}{
				"y": 1,
			},
		},
		tgt: map[string]interface{}{},
		want: map[string]interface{}{
			"x": map[string]interface{}{
				"y": 1,
			},
		},
	}, {
		name:    "identical nested slice",
		changed: false,
		src: map[string]interface{}{
			"x": map[string]interface{}{
				"y": []interface{}{"1", "2"},
			},
		},
		tgt: map[string]interface{}{
			"x": map[string]interface{}{
				"y": []interface{}{"1", "2"},
			},
		},
		want: map[string]interface{}{
			"x": map[string]interface{}{
				"y": []interface{}{"1", "2"},
			},
		},
	}, {
		name:    "add nested slice entry",
		changed: true,
		src: map[string]interface{}{
			"x": map[string]interface{}{
				"y": []interface{}{"1", "2"},
			},
		},
		tgt: map[string]interface{}{
			"x": map[string]interface{}{
				"y": []interface{}{"1"},
			},
		},
		want: map[string]interface{}{
			"x": map[string]interface{}{
				"y": []interface{}{"1", "2"},
			},
		},
	}, {
		name:    "update nested slice entry",
		changed: true,
		src: map[string]interface{}{
			"x": map[string]interface{}{
				"y": []interface{}{"1", "2", "3"},
			},
		},
		tgt: map[string]interface{}{
			"x": map[string]interface{}{
				"y": []interface{}{"3", "6", "9"},
			},
		},
		want: map[string]interface{}{
			"x": map[string]interface{}{
				"y": []interface{}{"1", "2", "3"},
			},
		},
	}, {
		name:    "add missing slice entry",
		changed: true,
		src: map[string]interface{}{
			"x": map[string]interface{}{
				"y": []interface{}{"1", "2"},
			},
		},
		tgt: map[string]interface{}{
			"x": map[string]interface{}{
				"x": 2,
			},
		},
		want: map[string]interface{}{
			"x": map[string]interface{}{
				"y": []interface{}{"1", "2"},
				"x": 2,
			},
		},
	}, {
		name:    "change map within list",
		changed: true,
		src: map[string]interface{}{
			"x": map[string]interface{}{
				"y": []interface{}{map[string]interface{}{"foo": "bar"}},
			},
		},
		tgt: map[string]interface{}{
			"x": map[string]interface{}{
				"y": []interface{}{map[string]interface{}{"foo": "baz", "one": 1}},
			},
		},
		want: map[string]interface{}{
			"x": map[string]interface{}{
				"y": []interface{}{map[string]interface{}{"foo": "bar"}},
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := fmt.Sprintf("%+v", test.tgt)
			actual := UpdateChanged(test.src, test.tgt)

			if actual != test.changed {
				t.Errorf("updateChanged() = %v, want: %v", actual, test.changed)
			}

			if !reflect.DeepEqual(test.tgt, test.want) {
				t.Errorf("from %+v to %s => %+v; wanted %+v", test.src, original, test.tgt, test.want)
			}
		})
	}
}
