package openapi

import (
	"reflect"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestOpenAPIUnmarshalYAML(t *testing.T) {
	yml := `---
openapi: 3.0.2
info:
  title: openapi spec
  version: 1.0.0
paths:
  /:
    get:
      responses:
        '200':
          description: ok`
	var got OpenAPI
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Error(err)
		return
	}
	want := OpenAPI{
		openapi: "3.0.2",
		info: &Info{
			title:   "openapi spec",
			version: "1.0.0",
		},
		paths: &Paths{
			paths: map[string]*PathItem{
				"/": &PathItem{
					get: &Operation{
						responses: &Responses{
							responses: map[string]*Response{
								"200": &Response{
									description: "ok",
								},
							},
						},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\n  got:  %#v\n  want: %#v", got, want)
		t.Log(got.paths.paths["/"].get.responses.responses, want.paths.paths["/"].get.responses.responses)
		t.Log(reflect.DeepEqual(got.paths.paths["/"].get.responses.responses, want.paths.paths["/"].get.responses.responses))
		return
	}
}
