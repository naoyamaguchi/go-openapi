package openapi

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestInfoUnmarshalYAML(t *testing.T) {
	yml := `title: Sample Pet Store App
description: This is a sample server for a pet store.
termsOfService: http://example.com/terms/
contact:
  name: API Support
  url: http://www.example.com/support
  email: support@example.com
license:
  name: Apache 2.0
  url: https://www.apache.org/licenses/LICENSE-2.0.html
version: 1.0.1`

	var info Info
	if err := yaml.Unmarshal([]byte(yml), &info); err != nil {
		t.Fatal(err)
	}

	if info.title != "Sample Pet Store App" {
		t.Errorf("unexpected info.title: %s", info.title)
		return
	}
	if info.description != "This is a sample server for a pet store." {
		t.Errorf("unexpected info.description: %s", info.description)
		return
	}
	if info.termsOfService != "http://example.com/terms/" {
		t.Errorf("unexpected info.termsOfService: %s", info.termsOfService)
		return
	}
	if info.contact.name != "API Support" {
		t.Errorf("unexpected info.contact.name: %s", info.contact.name)
		return
	}
	if info.contact.url != "http://www.example.com/support" {
		t.Errorf("unexpected info.contact.url: %s", info.contact.url)
		return
	}
	if info.contact.email != "support@example.com" {
		t.Errorf("unexpected info.contact.email: %s", info.contact.email)
		return
	}
	if info.license.name != "Apache 2.0" {
		t.Errorf("unexpected info.license.name: %s", info.license.name)
		return
	}
	if info.license.url != "https://www.apache.org/licenses/LICENSE-2.0.html" {
		t.Errorf("unexpected info.license.url: %s", info.license.url)
		return
	}
	if info.version != "1.0.1" {
		t.Errorf("unexpected info.version: %s", info.version)
		return
	}
}

func TestContactUnmarshalYAML(t *testing.T) {
	yml := `name: API Support
url: http://www.example.com/support
email: support@example.com`

	var contact Contact
	if err := yaml.Unmarshal([]byte(yml), &contact); err != nil {
		t.Fatal(err)
	}

	if contact.name != "API Support" {
		t.Errorf("unexpected contact.name: %s", contact.name)
		return
	}
	if contact.url != "http://www.example.com/support" {
		t.Errorf("unexpected contact.url: %s", contact.url)
		return
	}
	if contact.email != "support@example.com" {
		t.Errorf("unexpected contact.email: %s", contact.email)
		return
	}
}

func TestLicenseUnmarshalYAML(t *testing.T) {
	yml := `name: Apache 2.0
url: https://www.apache.org/licenses/LICENSE-2.0.html`

	var license License
	if err := yaml.Unmarshal([]byte(yml), &license); err != nil {
		t.Fatal(err)
	}

	if license.name != "Apache 2.0" {
		t.Errorf("unexpected license.name: %s", license.name)
		return
	}
	if license.url != "https://www.apache.org/licenses/LICENSE-2.0.html" {
		t.Errorf("unexpected license.url: %s", license.url)
		return
	}
}

func TestServerUnmarshalYAML(t *testing.T) {
	t.Run("single server", func(t *testing.T) {
		yml := `url: https://development.gigantic-server.com/v1
description: Development server`

		var server Server
		if err := yaml.Unmarshal([]byte(yml), &server); err != nil {
			t.Fatal(err)
		}

		if server.url != "https://development.gigantic-server.com/v1" {
			t.Errorf("unexpected server.url: %s", server.url)
			return
		}
		if server.description != "Development server" {
			t.Errorf("unexpected server.description: %s", server.description)
			return
		}
	})
	t.Run("servers", func(t *testing.T) {
		yml := `servers:
- url: https://development.gigantic-server.com/v1
  description: Development server
- url: https://staging.gigantic-server.com/v1
  description: Staging server
- url: https://api.gigantic-server.com/v1
  description: Production server`

		var target struct {
			Servers []*Server
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		servers := target.Servers
		t.Run("0", func(t *testing.T) {
			server := servers[0]
			if server.url != "https://development.gigantic-server.com/v1" {
				t.Errorf("unexpected server.url: %s", server.url)
				return
			}
			if server.description != "Development server" {
				t.Errorf("unexpected server.description: %s", server.description)
				return
			}
		})
		t.Run("1", func(t *testing.T) {
			server := servers[1]
			if server.url != "https://staging.gigantic-server.com/v1" {
				t.Errorf("unexpected server.url: %s", server.url)
				return
			}
			if server.description != "Staging server" {
				t.Errorf("unexpected server.description: %s", server.description)
				return
			}
		})
		t.Run("2", func(t *testing.T) {
			server := servers[2]
			if server.url != "https://api.gigantic-server.com/v1" {
				t.Errorf("unexpected server.url: %s", server.url)
				return
			}
			if server.description != "Production server" {
				t.Errorf("unexpected server.description: %s", server.description)
				return
			}
		})
	})
	t.Run("with variables", func(t *testing.T) {
		yml := `servers:
- url: https://{username}.gigantic-server.com:{port}/{basePath}
  description: The production API server
  variables:
    username:
      # note! no enum here means it is an open value
      default: demo
      description: this value is assigned by the service provider, in this example "gigantic-server.com"
    port:
      enum:
        - '8443'
        - '443'
      default: '8443'
    basePath:
      # open meaning there is the opportunity to use special base paths as assigned by the provider, default is "v2"
      default: v2`

		var target struct {
			Servers []*Server
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}

		server := target.Servers[0]
		if server.url != "https://{username}.gigantic-server.com:{port}/{basePath}" {
			t.Errorf("unexpected server url: %s", server.url)
			return
		}
		if server.description != "The production API server" {
			t.Errorf("unexpected server.descripion: %s", server.description)
			return
		}
		if server.variables["username"].default_ != "demo" {
			t.Errorf("unexpected server.variables.username.default: %s", server.variables["username"].default_)
			return
		}
		if server.variables["username"].description != `this value is assigned by the service provider, in this example "gigantic-server.com"` {
			t.Errorf("unexpected server.variables.username.description: %s", server.variables["username"].description)
			return
		}
		if len(server.variables["port"].enum) != 2 {
			t.Errorf("unexpected length of server.variables.port.enum: %d", len(server.variables["port"].enum))
			return
		}
		if !reflect.DeepEqual(server.variables["port"].enum, []string{"8443", "443"}) {
			t.Errorf("unexpected server.variables.port.enum: %#v", server.variables["port"].enum)
			return
		}
		if server.variables["port"].default_ != "8443" {
			t.Errorf("unexpected server.variables.port.default: %s", server.variables["port"].default_)
			return
		}
		if server.variables["basePath"].default_ != "v2" {
			t.Errorf("unexpected server.variables.basepath.default: %s", server.variables["basePath"].default_)
			return
		}
	})
}

func TestComponentsUnmarshalYAML(t *testing.T) {
	yml := `components:
  schemas:
    GeneralError:
      type: object
      properties:
        code:
          type: integer
          format: int32
        message:
          type: string
    Category:
      type: object
      properties:
        id:
          type: integer
          format: int64
        name:
          type: string
    Tag:
      type: object
      properties:
        id:
          type: integer
          format: int64
        name:
          type: string
  parameters:
    skipParam:
      name: skip
      in: query
      description: number of items to skip
      required: true
      schema:
        type: integer
        format: int32
    limitParam:
      name: limit
      in: query
      description: max records to return
      required: true
      schema:
        type: integer
        format: int32
  responses:
    NotFound:
      description: Entity not found.
    IllegalInput:
      description: Illegal input for operation.
    GeneralError:
      description: General Error
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/GeneralError'
  securitySchemes:
    api_key:
      type: apiKey
      name: api_key
      in: header
    petstore_auth:
      type: oauth2
      flows:
        implicit:
          authorizationUrl: http://example.org/api/oauth/dialog
          scopes:
            write:pets: modify pets in your account
            read:pets: read your pets`
	var target struct {
		Components *Components
	}
	if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
		t.Fatal(err)
	}
	components := target.Components
	schemas := components.schemas
	t.Run("schemas.GeneralError", func(t *testing.T) {
		generalError, ok := schemas["GeneralError"]
		if !ok {
			t.Error("schemas.GeneralError is not found")
			return
		}
		if generalError.type_ != "object" {
			t.Errorf("unexpected schema.type: %s", generalError.type_)
			return
		}
		code, ok := generalError.properties["code"]
		if !ok {
			t.Error("schemas.GeneralError.properties.code is not found")
			return
		}
		if code.type_ != "integer" {
			t.Errorf("unexpected schemas.GeneralError.properties.code.type: %s", code.type_)
			return
		}
		if code.format != "int32" {
			t.Errorf("unexpected schemas.GeneralError.properties.code.format: %s", code.format)
			return
		}
		message, ok := generalError.properties["message"]
		if !ok {
			t.Error("schemas.GeneralError.properties.message is not found")
			return
		}
		if message.type_ != "string" {
			t.Errorf("unexpected schemas.GeneralError.properties.message.type: %s", message.type_)
			return
		}
	})
	t.Run("schemas.Category", func(t *testing.T) {
		category, ok := schemas["Category"]
		if !ok {
			t.Error("schemas.Category is not found")
			return
		}
		if category.type_ != "object" {
			t.Errorf("unexpected schema.Category.type: %s", category.type_)
			return
		}
		id, ok := category.properties["id"]
		if !ok {
			t.Error("schemas.Category.properties.id is not found")
			return
		}
		if id.type_ != "integer" {
			t.Errorf("unexpected schemas.Category.properties.id.type: %s", id.type_)
			return
		}
		if id.format != "int64" {
			t.Errorf("unexpected schemas.Category.properties.id.format: %s", id.format)
			return
		}
		name, ok := category.properties["name"]
		if !ok {
			t.Error("schemas.Category.properties.name is not found")
			return
		}
		if name.type_ != "string" {
			t.Errorf("unexpected schemas.Category.properties.name.type: %s", name.type_)
			return
		}
	})
	parameters := components.parameters
	t.Run("parameters.skipParam", func(t *testing.T) {
		skipParam, ok := parameters["skipParam"]
		if !ok {
			t.Error("parameters.skipParam is not found")
			return
		}
		if skipParam.name != "skip" {
			t.Errorf("unexpected parameters.skipParam.name: %s", skipParam.name)
			return
		}
		if skipParam.in != "query" {
			t.Errorf("unexpected parameters.skipParam.in: %s", skipParam.in)
			return
		}
		if skipParam.description != "number of items to skip" {
			t.Errorf("unexpected parameters.skipParam.description: %s", skipParam.description)
			return
		}
		if skipParam.required != true {
			t.Errorf("unexpected parameters.skipParam.required: %t", skipParam.required)
			return
		}
		schema := skipParam.schema
		if schema.type_ != "integer" {
			t.Errorf("unexpected parameters.skipParam.schema.type: %s", schema.type_)
			return
		}
		if schema.format != "int32" {
			t.Errorf("unexpected parameters.skipParam.schema.format: %s", schema.format)
			return
		}
	})
	t.Run("parameters.limitParam", func(t *testing.T) {
		limitParam, ok := parameters["limitParam"]
		if !ok {
			t.Error("parameters.limitParam is not found")
			return
		}
		if limitParam.name != "limit" {
			t.Errorf("unexpected parameters.limitParam.name: %s", limitParam.name)
			return
		}
		if limitParam.in != "query" {
			t.Errorf("unexpected parameters.limitParam.in: %s", limitParam.in)
			return
		}
		if limitParam.description != "max records to return" {
			t.Errorf("unexpected parameters.limitParam.description: %s", limitParam.description)
			return
		}
		if limitParam.required != true {
			t.Errorf("unexpected parameters.limitParam.required: %t", limitParam.required)
			return
		}
		schema := limitParam.schema
		if schema.type_ != "integer" {
			t.Errorf("unexpected parameters.limitParam.schema.type: %s", schema.type_)
			return
		}
		if schema.format != "int32" {
			t.Errorf("unexpected parameters.limitParam.schema.format: %s", schema.format)
			return
		}
	})
	responses := components.responses
	t.Run("responses.NotFound", func(t *testing.T) {
		notFound, ok := responses["NotFound"]
		if !ok {
			t.Error("responses.NotFound is not found")
			return
		}
		if notFound.description != "Entity not found." {
			t.Errorf("unexpected responses.NotFound.description: %s", notFound.description)
			return
		}
	})
	t.Run("responses.IllegalInput", func(t *testing.T) {
		illegalInput, ok := responses["IllegalInput"]
		if !ok {
			t.Error("responses.IllegalInput is not found")
			return
		}
		if illegalInput.description != "Illegal input for operation." {
			t.Errorf("unexpected responses.IllegalInput.description: %s", illegalInput.description)
			return
		}
	})
	t.Run("responses.GeneralError", func(t *testing.T) {
		generalError, ok := responses["GeneralError"]
		if !ok {
			t.Error("responses.GeneralError is not found")
			return
		}
		if generalError.description != "General Error" {
			t.Errorf("unexpected responses.GeneralError.description: %s", generalError.description)
			return
		}
		mediaType, ok := generalError.content["application/json"]
		if !ok {
			t.Error("responses.GeneralError.content.application/json is not found")
			return
		}
		if mediaType.schema.reference != "#/components/schemas/GeneralError" {
			t.Errorf("unexpected responses.GeneralError.content.application/json.schema.$ref")
			return
		}
	})
	securitySchemes := components.securitySchemes
	t.Run("securitySchemes.api_key", func(t *testing.T) {
		apiKey, ok := securitySchemes["api_key"]
		if !ok {
			t.Error("securitySchemes.api_key is not found")
			return
		}
		if apiKey.type_ != "apiKey" {
			t.Errorf("unexpected securitySchemes.api_key.type: %s", apiKey.type_)
			return
		}
		if apiKey.name != "api_key" {
			t.Errorf("unexpected securitySchemes.api_key.name: %s", apiKey.name)
			return
		}
		if apiKey.in != "header" {
			t.Errorf("unexpected securitySchemes.api_key.in: %s", apiKey.in)
			return
		}
	})
	t.Run("securitySchemes.petstore_auth", func(t *testing.T) {
		petstoreAuth, ok := securitySchemes["petstore_auth"]
		if !ok {
			t.Error("securitySchemes.petstore_auth is not found")
			return
		}
		if petstoreAuth.type_ != "oauth2" {
			t.Errorf("unexpected securitySchemes.petstore_auth.type: %s", petstoreAuth.type_)
			return
		}
		if petstoreAuth.flows.implicit.authorizationURL != "http://example.org/api/oauth/dialog" {
			t.Errorf("unexpected securitySchemes.petstore_auth.flows.implicit.authorizationURL: %s", petstoreAuth.flows.implicit.authorizationURL)
			return
		}
		scopes := petstoreAuth.flows.implicit.scopes
		write, ok := scopes["write:pets"]
		if !ok {
			t.Error("securitySchemes.petstore_auth.flows.implicit.scopes.write:pets is not found")
			return
		}
		if write != "modify pets in your account" {
			t.Errorf("unexpected securitySchemes.petstore_auth.flows.implicit.scopes.write:pets: %s", write)
		}
		read, ok := scopes["read:pets"]
		if !ok {
			t.Error("securitySchemes.petstore_auth.flows.implicit.scopes.read:pets is not found")
			return
		}
		if read != "read your pets" {
			t.Errorf("unexpected securitySchemes.petstore_auth.flows.implicit.scopes.read:pets: %s", write)
		}
	})
}

func TestPathsUnmarshalYAML(t *testing.T) {
	yml := `/pets:
  get:
    description: Returns all pets from the system that the user has access to
    responses:
      '200':
        description: A list of pets.
        content:
          application/json:
            schema:
              type: array
              items:
                $ref: '#/components/schemas/pet'`
	var paths Paths
	if err := yaml.Unmarshal([]byte(yml), &paths); err != nil {
		t.Fatal(err)
	}
	if _, ok := paths.paths["/pets"]; !ok {
		t.Error("paths./pets is not found")
		return
	}
	op := paths.paths["/pets"].get
	if op.description != "Returns all pets from the system that the user has access to" {
		t.Errorf("unexpected paths./pets.get.description: %s", op.description)
		return
	}
	response, ok := op.responses.responses["200"]
	if !ok {
		t.Error("paths./pets.get.responses.200 is not found")
		return
	}
	if _, ok := response.content["application/json"]; !ok {
		t.Error("paths./pets.get.responses.200.content.application/json is not found")
		return
	}
	schema := response.content["application/json"].schema
	if schema.type_ != "array" {
		t.Errorf("unexpected paths./pets.get.responses.200.content.application/json.schema.type: %s", schema.type_)
		return
	}
	if schema.items.reference != "#/components/schemas/pet" {
		t.Errorf("unexpected paths./pets.get.responses.200.content.application/json.schema.items.$ref: %s", schema.reference)
		return
	}
}

func TestIsOneOf(t *testing.T) {
	tests := []struct {
		s    string
		list []string
		want bool
	}{
		{
			s:    "",
			list: []string{},
			want: false,
		},
		{
			s:    "a",
			list: []string{"a", "b"},
			want: true,
		},
		{
			s:    "c",
			list: []string{"a", "b"},
			want: false,
		},
		{
			s:    "a",
			list: nil,
			want: false,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := isOneOf(tt.s, tt.list)
			if got != tt.want {
				t.Errorf("unexpected: %t != %t", got, tt.want)
				return
			}
		})
	}
}

func TestMatchRuntimerExpr(t *testing.T) {
	tests := []struct {
		expr string
		want bool
	}{
		{
			expr: "$method",
			want: true,
		},
		{
			expr: "$request.header.accept",
			want: true,
		},
		{
			expr: "$request.path.id",
			want: true,
		},
		{
			expr: "$request.body#/user/uuid",
			want: true,
		},
		{
			expr: "$url",
			want: true,
		},
		{
			expr: "$response.body#/status",
			want: true,
		},
		{
			expr: "$response.header.Server",
			want: true,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"/"+tt.expr, func(t *testing.T) {
			got := matchRuntimeExpr(tt.expr)
			if got != tt.want {
				t.Errorf("unexpected: %t != %t", got, tt.want)
				return
			}
		})
	}
}
