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

func TestUsptoUnmarshalYAML(t *testing.T) {
	b, err := ioutil.ReadFile("test/testdata/uspto.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var got OpenAPI
	if err := yaml.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}

	want := OpenAPI{
		openapi: "3.0.1",
		servers: []*Server{
			{
				url: "{scheme}://developer.uspto.gov/ds-api",
				variables: map[string]*ServerVariable{
					"scheme": {
						description: "The Data Set API is accessible via https and http",
						enum: []string{
							"https",
							"http",
						},
						default_: "https",
					},
				},
			},
		},
		info: &Info{
			description: `The Data Set API (DSAPI) allows the public users to discover and search
USPTO exported data sets. This is a generic API that allows USPTO users to
make any CSV based data files searchable through API. With the help of GET
call, it returns the list of data fields that are searchable. With the help
of POST call, data can be fetched based on the filters on the field names.
Please note that POST call is used to search the actual data. The reason for
the POST call is that it allows users to specify any complex search criteria
without worry about the GET size limitations as well as encoding of the
input parameters.`,
			version: "1.0.0",
			title:   "USPTO Data Set API",
			contact: &Contact{
				name:  "Open Data Portal",
				url:   "https://developer.uspto.gov",
				email: "developer@uspto.gov",
			},
		},
		tags: []*Tag{
			{
				name:        "metadata",
				description: "Find out about the data sets",
			},
			{
				name:        "search",
				description: "Search a data set",
			},
		},
		paths: &Paths{
			paths: map[string]*PathItem{
				"/": {
					get: &Operation{
						tags:        []string{"metadata"},
						operationID: "list-data-sets",
						summary:     "List available data sets",
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "Returns a list of data sets",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												reference: "#/components/schemas/dataSetList",
											},
											example: map[string]interface{}{
												"total": uint64(2),
												"apis": []interface{}{
													map[string]interface{}{
														"apiKey":              "oa_citations",
														"apiVersionNumber":    "v1",
														"apiUrl":              "https://developer.uspto.gov/ds-api/oa_citations/v1/fields",
														"apiDocumentationUrl": "https://developer.uspto.gov/ds-api-docs/index.html?url=https://developer.uspto.gov/ds-api/swagger/docs/oa_citations.json",
													},
													map[string]interface{}{
														"apiKey":              "cancer_moonshot",
														"apiVersionNumber":    "v1",
														"apiUrl":              "https://developer.uspto.gov/ds-api/cancer_moonshot/v1/fields",
														"apiDocumentationUrl": "https://developer.uspto.gov/ds-api-docs/index.html?url=https://developer.uspto.gov/ds-api/swagger/docs/cancer_moonshot.json",
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
				"/{dataset}/{version}/fields": {
					get: &Operation{
						tags: []string{"metadata"},
						summary: `Provides the general information about the API and the list of fields
that can be used to query the dataset.`,
						description: `This GET API returns the list of all the searchable field names that are
in the oa_citations. Please see the 'fields' attribute which returns an
array of field names. Each field or a combination of fields can be
searched using the syntax options shown below.`,
						operationID: "list-searchable-fields",
						parameters: []*Parameter{
							{
								name:        "dataset",
								in:          "path",
								description: "Name of the dataset.",
								required:    true,
								example:     "oa_citations",
								schema: &Schema{
									type_: "string",
								},
							},
							{
								name:        "version",
								in:          "path",
								description: "Version of the dataset.",
								required:    true,
								example:     "v1",
								schema: &Schema{
									type_: "string",
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: `The dataset API for the given version is found and it is accessible
to consume.`,
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												type_: "string",
											},
										},
									},
								},
								"404": {
									description: `The combination of dataset name and version is not found in the
system or it is not published yet to be consumed by public.`,
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												type_: "string",
											},
										},
									},
								},
							},
						},
					},
				},
				"/{dataset}/{version}/records": {
					post: &Operation{
						tags: []string{"search"},
						summary: `Provides search capability for the data set with the given search
criteria.`,
						description: `This API is based on Solr/Lucense Search. The data is indexed using
SOLR. This GET API returns the list of all the searchable field names
that are in the Solr Index. Please see the 'fields' attribute which
returns an array of field names. Each field or a combination of fields
can be searched using the Solr/Lucene Syntax. Please refer
https://lucene.apache.org/core/3_6_2/queryparsersyntax.html#Overview for
the query syntax. List of field names that are searchable can be
determined using above GET api.`,
						operationID: "perform-search",
						parameters: []*Parameter{
							{
								name:        "version",
								in:          "path",
								description: "Version of the dataset.",
								required:    true,
								schema: &Schema{
									type_:    "string",
									default_: "v1",
								},
							},
							{
								name:        "dataset",
								in:          "path",
								description: "Name of the dataset. In this case, the default value is oa_citations",
								required:    true,
								schema: &Schema{
									type_:    "string",
									default_: "oa_citations",
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "successful operation",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												type_: "array",
												items: &Schema{
													type_: "object",
													additionalProperties: &Schema{
														type_: "object",
													},
												},
											},
										},
									},
								},
								"404": {
									description: "No matching record found for the given criteria.",
								},
							},
						},
						requestBody: &RequestBody{
							content: map[string]*MediaType{
								"application/x-www-form-urlencoded": {
									schema: &Schema{
										type_: "object",
										properties: map[string]*Schema{

											"criteria": {
												description: `Uses Lucene Query Syntax in the format of
propertyName:value, propertyName:[num1 TO num2] and date
range format: propertyName:[yyyyMMdd TO yyyyMMdd]. In the
response please see the 'docs' element which has the list of
record objects. Each record structure would consist of all
the fields and their corresponding values.`,
												type_:    "string",
												default_: "*:*",
											},
											"start": {
												description: "Starting record number. Default value is 0.",
												type_:       "integer",
												default_:    "0",
											},
											"rows": {
												description: `Specify number of rows to be returned. If you run the search
with default values, in the response you will see 'numFound'
attribute which will tell the number of records available in
the dataset.`,
												type_:    "integer",
												default_: "100",
											},
										},
										required: []string{"criteria"},
									},
								},
							},
						},
					},
				},
			},
		},
		components: &Components{
			schemas: map[string]*Schema{
				"dataSetList": {
					type_: "object",
					properties: map[string]*Schema{
						"total": {
							type_: "integer",
						},
						"apis": {
							type_: "array",
							items: &Schema{
								type_: "object",
								properties: map[string]*Schema{
									"apiKey": {
										type_:       "string",
										description: "To be used as a dataset parameter value",
									},
									"apiVersionNumber": {
										type_:       "string",
										description: "To be used as a version parameter value",
									},
									"apiUrl": {
										type_:       "string",
										format:      "uriref",
										description: "The URL describing the dataset's fields",
									},
									"apiDocumentationUrl": {
										type_:       "string",
										format:      "uriref",
										description: "A URL to the API console for each API",
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
