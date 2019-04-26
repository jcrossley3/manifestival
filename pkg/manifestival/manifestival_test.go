package manifestival_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/jcrossley3/manifestival/pkg/manifestival"
)

func TestFinding(t *testing.T) {
	f, err := NewYamlManifest("testdata/", true, nil)
	if err != nil {
		t.Errorf("NewYamlManifest() = %v, wanted no error", err)
	}

	f.Filter(ByNamespace("fubar"))
	actual := f.Find("v1", "A", "foo")
	if actual == nil {
		t.Error("Failed to find resource")
	}
	if actual.GetNamespace() != "fubar" {
		t.Errorf("Resource has wrong namespace: %s", actual)
	}
	if f.Find("NO", "NO", "NO") != nil {
		t.Error("Missing resource found")
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
