package openapi

import (
	"io/ioutil"
	"reflect"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestCallbackExampleUnmarshalYAML(t *testing.T) {
	b, err := ioutil.ReadFile("test/testdata/callback-example.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var got OpenAPI
	if err := yaml.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}

	want := OpenAPI{
		openapi: "3.0.0",
		info: &Info{
			title:   "Callback Example",
			version: "1.0.0",
		},
		paths: &Paths{
			paths: map[string]*PathItem{
				"/streams": {
					post: &Operation{
						description: "subscribes a client to receive out-of-band data",
						parameters: []*Parameter{
							{
								name:        "callbackUrl",
								in:          "query",
								required:    true,
								description: "the location where data will be sent.  Must be network accessible\nby the source server",
								schema: &Schema{
									type_:   "string",
									format:  "uri",
									example: "https://tonys-server.com",
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"201": {
									description: "subscription successfully created",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												description: "subscription information",
												required: []string{
													"subscriptionId",
												},
												properties: map[string]*Schema{
													"subscriptionId": {
														description: "this unique identifier allows management of the subscription",
														type_:       "string",
														example:     "2531329f-fb09-4ef7-887e-84e648214436",
													},
												},
											},
										},
									},
								},
							},
						},
						callbacks: map[string]*Callback{
							"onData": {
								callback: map[string]*PathItem{
									"{$request.query.callbackUrl}/data": {
										post: &Operation{
											requestBody: &RequestBody{
												description: "subscription payload",
												content: map[string]*MediaType{
													"application/json": {
														schema: &Schema{
															type_: "object",
															properties: map[string]*Schema{
																"timestamp": {
																	type_:  "string",
																	format: "date-time",
																},
																"userData": {
																	type_: "string",
																},
															},
														},
													},
												},
											},
											responses: &Responses{
												responses: map[string]*Response{
													"202": {
														description: "Your server implementation should return this HTTP status code\nif the data was received successfully",
													},
													"204": {
														description: "Your server should return this HTTP status code if no longer interested\nin further updates",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected:\n  got:  %#v\n  want: %#v", got, want)
		return
	}
}
