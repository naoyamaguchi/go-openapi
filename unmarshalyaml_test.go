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
			t.Error("components.schemas.GeneralError is not found")
			return
		}
		if generalError.type_ != "object" {
			t.Errorf("unexpected components.schema.type: %s", generalError.type_)
			return
		}
		code, ok := generalError.properties["code"]
		if !ok {
			t.Error("components.schemas.GeneralError.properties.code is not found")
			return
		}
		if code.type_ != "integer" {
			t.Errorf("unexpected components.schemas.GeneralError.properties.code.type: %s", code.type_)
			return
		}
		if code.format != "int32" {
			t.Errorf("unexpected components.schemas.GeneralError.properties.code.format: %s", code.format)
			return
		}
		message, ok := generalError.properties["message"]
		if !ok {
			t.Error("components.schemas.GeneralError.properties.message is not found")
			return
		}
		if message.type_ != "string" {
			t.Errorf("unexpected components.schemas.GeneralError.properties.message.type: %s", message.type_)
			return
		}
	})
	t.Run("schemas.Category", func(t *testing.T) {
		category, ok := schemas["Category"]
		if !ok {
			t.Error("components.schemas.Category is not found")
			return
		}
		if category.type_ != "object" {
			t.Errorf("unexpected components.schema.Category.type: %s", category.type_)
			return
		}
		id, ok := category.properties["id"]
		if !ok {
			t.Error("components.schemas.Category.properties.id is not found")
			return
		}
		if id.type_ != "integer" {
			t.Errorf("unexpected components.schemas.Category.properties.id.type: %s", id.type_)
			return
		}
		if id.format != "int64" {
			t.Errorf("unexpected components.schemas.Category.properties.id.format: %s", id.format)
			return
		}
		name, ok := category.properties["name"]
		if !ok {
			t.Error("components.schemas.Category.properties.name is not found")
			return
		}
		if name.type_ != "string" {
			t.Errorf("unexpected components.schemas.Category.properties.name.type: %s", name.type_)
			return
		}
	})
	parameters := components.parameters
	t.Run("parameters.skipParam", func(t *testing.T) {
		skipParam, ok := parameters["skipParam"]
		if !ok {
			t.Error("components.parameters.skipParam is not found")
			return
		}
		if skipParam.name != "skip" {
			t.Errorf("unexpected components.parameters.skipParam.name: %s", skipParam.name)
			return
		}
		if skipParam.in != "query" {
			t.Errorf("unexpected components.parameters.skipParam.in: %s", skipParam.in)
			return
		}
		if skipParam.description != "number of items to skip" {
			t.Errorf("unexpected components.parameters.skipParam.description: %s", skipParam.description)
			return
		}
		if skipParam.required != true {
			t.Errorf("unexpected components.parameters.skipParam.required: %t", skipParam.required)
			return
		}
		schema := skipParam.schema
		if schema.type_ != "integer" {
			t.Errorf("unexpected components.parameters.skipParam.schema.type: %s", schema.type_)
			return
		}
		if schema.format != "int32" {
			t.Errorf("unexpected components.parameters.skipParam.schema.format: %s", schema.format)
			return
		}
	})
	t.Run("parameters.limitParam", func(t *testing.T) {
		limitParam, ok := parameters["limitParam"]
		if !ok {
			t.Error("components.parameters.limitParam is not found")
			return
		}
		if limitParam.name != "limit" {
			t.Errorf("unexpected components.parameters.limitParam.name: %s", limitParam.name)
			return
		}
		if limitParam.in != "query" {
			t.Errorf("unexpected components.parameters.limitParam.in: %s", limitParam.in)
			return
		}
		if limitParam.description != "max records to return" {
			t.Errorf("unexpected components.parameters.limitParam.description: %s", limitParam.description)
			return
		}
		if limitParam.required != true {
			t.Errorf("unexpected components.parameters.limitParam.required: %t", limitParam.required)
			return
		}
		schema := limitParam.schema
		if schema.type_ != "integer" {
			t.Errorf("unexpected components.parameters.limitParam.schema.type: %s", schema.type_)
			return
		}
		if schema.format != "int32" {
			t.Errorf("unexpected components.parameters.limitParam.schema.format: %s", schema.format)
			return
		}
	})
	responses := components.responses
	t.Run("responses.NotFound", func(t *testing.T) {
		notFound, ok := responses["NotFound"]
		if !ok {
			t.Error("components.responses.NotFound is not found")
			return
		}
		if notFound.description != "Entity not found." {
			t.Errorf("unexpected components.responses.NotFound.description: %s", notFound.description)
			return
		}
	})
	t.Run("responses.IllegalInput", func(t *testing.T) {
		illegalInput, ok := responses["IllegalInput"]
		if !ok {
			t.Error("components.responses.IllegalInput is not found")
			return
		}
		if illegalInput.description != "Illegal input for operation." {
			t.Errorf("unexpected components.responses.IllegalInput.description: %s", illegalInput.description)
			return
		}
	})
	t.Run("responses.GeneralError", func(t *testing.T) {
		generalError, ok := responses["GeneralError"]
		if !ok {
			t.Error("components.responses.GeneralError is not found")
			return
		}
		if generalError.description != "General Error" {
			t.Errorf("unexpected components.responses.GeneralError.description: %s", generalError.description)
			return
		}
		mediaType, ok := generalError.content["application/json"]
		if !ok {
			t.Error("components.responses.GeneralError.content.application/json is not found")
			return
		}
		if mediaType.schema.reference != "#/components/schemas/GeneralError" {
			t.Errorf("unexpected components.responses.GeneralError.content.application/json.schema.$ref")
			return
		}
	})
	securitySchemes := components.securitySchemes
	t.Run("securitySchemes.api_key", func(t *testing.T) {
		apiKey, ok := securitySchemes["api_key"]
		if !ok {
			t.Error("components.securitySchemes.api_key is not found")
			return
		}
		if apiKey.type_ != "apiKey" {
			t.Errorf("unexpected components.securitySchemes.api_key.type: %s", apiKey.type_)
			return
		}
		if apiKey.name != "api_key" {
			t.Errorf("unexpected components.securitySchemes.api_key.name: %s", apiKey.name)
			return
		}
		if apiKey.in != "header" {
			t.Errorf("unexpected components.securitySchemes.api_key.in: %s", apiKey.in)
			return
		}
	})
	t.Run("securitySchemes.petstore_auth", func(t *testing.T) {
		petstoreAuth, ok := securitySchemes["petstore_auth"]
		if !ok {
			t.Error("components.securitySchemes.petstore_auth is not found")
			return
		}
		if petstoreAuth.type_ != "oauth2" {
			t.Errorf("unexpected components.securitySchemes.petstore_auth.type: %s", petstoreAuth.type_)
			return
		}
		if petstoreAuth.flows.implicit.authorizationURL != "http://example.org/api/oauth/dialog" {
			t.Errorf("unexpected components.securitySchemes.petstore_auth.flows.implicit.authorizationURL: %s", petstoreAuth.flows.implicit.authorizationURL)
			return
		}
		scopes := petstoreAuth.flows.implicit.scopes
		write, ok := scopes["write:pets"]
		if !ok {
			t.Error("components.securitySchemes.petstore_auth.flows.implicit.scopes.write:pets is not found")
			return
		}
		if write != "modify pets in your account" {
			t.Errorf("unexpected components.securitySchemes.petstore_auth.flows.implicit.scopes.write:pets: %s", write)
		}
		read, ok := scopes["read:pets"]
		if !ok {
			t.Error("components.securitySchemes.petstore_auth.flows.implicit.scopes.read:pets is not found")
			return
		}
		if read != "read your pets" {
			t.Errorf("unexpected components.securitySchemes.petstore_auth.flows.implicit.scopes.read:pets: %s", write)
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

func TestPathItemUnmarshalYAML(t *testing.T) {
	yml := `get:
  description: Returns pets based on ID
  summary: Find pets by ID
  operationId: getPetsById
  responses:
    '200':
      description: pet response
      content:
        '*/*' :
          schema:
            type: array
            items:
              $ref: '#/components/schemas/Pet'
    default:
      description: error payload
      content:
        'text/html':
          schema:
            $ref: '#/components/schemas/ErrorModel'
parameters:
- name: id
  in: path
  description: ID of pet to use
  required: true
  schema:
    type: array
    # This is in example but maybe mistake
    # style: simple
    items:
      type: string  `
	var pathItem PathItem
	if err := yaml.Unmarshal([]byte(yml), &pathItem); err != nil {
		t.Fatal(err)
	}
	t.Run("get", func(t *testing.T) {
		operation := pathItem.get
		if operation.description != "Returns pets based on ID" {
			t.Errorf("unexpected pathItem.get.description: %s", operation.description)
			return
		}
		if operation.summary != "Find pets by ID" {
			t.Errorf("unexpected pathItem.get.summary: %s", operation.summary)
			return
		}
		if operation.operationID != "getPetsById" {
			t.Errorf("unexpected pathItem.get.operationId: %s", operation.operationID)
			return
		}
		t.Run("200", func(t *testing.T) {
			response, ok := pathItem.get.responses.responses["200"]
			if !ok {
				t.Error("pathItem.get.responses.200 is not found")
				return
			}
			if response.description != "pet response" {
				t.Errorf("unexpected pathItem.get.responses.200.description: %s", response.description)
				return
			}
			if _, ok := response.content["*/*"]; !ok {
				t.Error("pathItem.get.responses.200.content.*/* is not found")
				return
			}
			schema := response.content["*/*"].schema
			if schema.type_ != "array" {
				t.Errorf("unexpected pathItem.get.responses.200.content.*/*.schema.type: %s", schema.type_)
				return
			}
			if schema.items.reference != "#/components/schemas/Pet" {
				t.Errorf("unexpected pathItem.get.responses.200.content.*/*.schema.items.$ref: %s", schema.items.reference)
				return
			}
		})
		t.Run("default", func(t *testing.T) {
			response, ok := pathItem.get.responses.responses["default"]
			if !ok {
				t.Error("pathItem.get.responses.default is not found")
				return
			}
			if response.description != "error payload" {
				t.Errorf("unexpected pathItem.get.responses.default.description: %s", response.description)
				return
			}
			if _, ok := response.content["text/html"]; !ok {
				t.Error("pathItem.get.responses.default.content.text/html is not found")
				return
			}
			if response.content["text/html"].schema.reference != "#/components/schemas/ErrorModel" {
				t.Errorf("unexpected pathItem.get.responses.default.content.text/html.schema.$ref: %s", response.content["text/html"].schema.reference)
				return
			}
		})
	})
	t.Run("parameters", func(t *testing.T) {
		parameters := pathItem.parameters
		id := parameters[0]
		if id.name != "id" {
			t.Errorf("unexpected pathItem.parameters.0.name: %s", id.name)
			return
		}
		if id.in != "path" {
			t.Errorf("unexpected pathItem.parameters.0.in: %s", id.in)
			return
		}
		if id.description != "ID of pet to use" {
			t.Errorf("unexpected pathItem.parameters.0.description: %s", id.description)
			return
		}
		if id.required != true {
			t.Errorf("unexpected pathItem.parameters.0.required: %t", id.required)
			return
		}
		if id.schema.type_ != "array" {
			t.Errorf("unexpected pathItem.parameters.0.schema.type: %s", id.schema.type_)
			return
		}
		if id.schema.items.type_ != "string" {
			t.Errorf("unexpected pathItem.parameters.0.schema.items.type: %s", id.schema.items.type_)
			return
		}
	})
}

func TestOperationUnmarshalYAML(t *testing.T) {
	yml := `tags:
- pet
summary: Updates a pet in the store with form data
operationId: updatePetWithForm
parameters:
- name: petId
  in: path
  description: ID of pet that needs to be updated
  required: true
  schema:
    type: string
requestBody:
  content:
    'application/x-www-form-urlencoded':
      schema:
       properties:
          name:
            description: Updated name of the pet
            type: string
          status:
            description: Updated status of the pet
            type: string
       required:
         - status
responses:
  '200':
    description: Pet updated.
    content:
      'application/json': {}
      'application/xml': {}
  '405':
    description: Method Not Allowed
    content:
      'application/json': {}
      'application/xml': {}
security:
- petstore_auth:
  - write:pets
  - read:pets`
	var operation Operation
	if err := yaml.Unmarshal([]byte(yml), &operation); err != nil {
		t.Fatal(err)
	}
	if operation.tags[0] != "pet" {
		t.Errorf("unexpected operation.tags.0: %s", operation.tags[0])
		return
	}
	if operation.summary != "Updates a pet in the store with form data" {
		t.Errorf("unexpected opration.summary: %s", operation.summary)
		return
	}
	if operation.operationID != "updatePetWithForm" {
		t.Errorf("unexpected operation.operationId: %s", operation.operationID)
		return
	}
	parameter := operation.parameters[0]
	if parameter.name != "petId" {
		t.Errorf("unexpected operation.parameters.0.name: %s", parameter.name)
		return
	}
	if parameter.in != "path" {
		t.Errorf("unexpected operation.parameters.0.in: %s", parameter.in)
		return
	}
	if parameter.description != "ID of pet that needs to be updated" {
		t.Errorf("unexpected operation.parameters.0.description: %s", parameter.description)
		return
	}
	if parameter.required != true {
		t.Errorf("unexpected operation.parameters.0.required: %t", parameter.required)
		return
	}
	if parameter.schema.type_ != "string" {
		t.Errorf("unexpected operation.parameters.0.schema.type: %s", parameter.schema.type_)
		return
	}
	if _, ok := operation.requestBody.content["application/x-www-form-urlencoded"]; !ok {
		t.Error("operation.requestBody.content.application/x-www-form-urlencoded is not found")
		return
	}
	schema := operation.requestBody.content["application/x-www-form-urlencoded"].schema
	name, ok := schema.properties["name"]
	if !ok {
		t.Error("operation.requestBody.content.application/x-www-form-urlencoded.schema.properties.name is not found")
		return
	}
	if name.description != "Updated name of the pet" {
		t.Errorf("unexpected operation.requestBody.content.application/x-www-form-urlencoded.schema.properties.name.description: %s", name.description)
		return
	}
	if name.type_ != "string" {
		t.Errorf("unexpected operation.requestBody.content.application/x-www-form-urlencoded.schema.properties.name.type: %s", name.type_)
		return
	}
	status, ok := schema.properties["status"]
	if !ok {
		t.Error("operation.requestBody.content.application/x-www-form-urlencoded.schema.properties.status is not found")
		return
	}
	if status.description != "Updated status of the pet" {
		t.Errorf("unexpected operation.requestBody.content.application/x-www-form-urlencoded.schema.properties.status.description: %s", status.description)
		return
	}
	if status.type_ != "string" {
		t.Errorf("unexpected operation.requestBody.content.application/x-www-form-urlencoded.schema.properties.status.type: %s", status.type_)
		return
	}
	if schema.required[0] != "status" {
		t.Errorf("unexpected operation.requestBody.content.application/x-www-form-urlencoded.schema.required.0: %s", schema.required[0])
		return
	}
	if _, ok := operation.responses.responses["200"]; !ok {
		t.Error("operation.responses.200 is not found")
		return
	}
	if operation.responses.responses["200"].description != "Pet updated." {
		t.Errorf("unexpected operation.responses.200.description: %s", operation.responses.responses["200"].description)
		return
	}
	if _, ok := operation.responses.responses["200"].content["application/json"]; !ok {
		t.Error("operation.responses.200.content.application/json is not found")
	}
	if _, ok := operation.responses.responses["200"].content["application/xml"]; !ok {
		t.Error("operation.responses.200.content.application/xml is not found")
		return
	}
	if _, ok := operation.responses.responses["405"]; !ok {
		t.Error("operation.responses.405 is not found")
		return
	}
	if operation.responses.responses["405"].description != "Method Not Allowed" {
		t.Errorf("unexpected operation.responses.405.description: %s", operation.responses.responses["405"].description)
		return
	}
	if _, ok := operation.responses.responses["405"].content["application/json"]; !ok {
		t.Error("operation.responses.405.content.application/json is not found")
	}
	if _, ok := operation.responses.responses["405"].content["application/xml"]; !ok {
		t.Error("operation.responses.405.content.application/xml is not found")
		return
	}
	securityRequirement, ok := operation.security[0].securityRequirement["petstore_auth"]
	if !ok {
		t.Error("operation.security.0.petstore_auth is not found")
		return
	}
	if securityRequirement[0] != "write:pets" {
		t.Errorf("unexpected operation.security.0.petstore_auth.0: %s", securityRequirement[0])
		return
	}
	if securityRequirement[1] != "read:pets" {
		t.Errorf("unexpected operation.security.0.petstore_auth.1: %s", securityRequirement[1])
		return
	}
}

func TestExternalDocumentationUnmarshalYAML(t *testing.T) {
	yml := `description: Find more info here
url: https://example.com`
	var externalDocumentation ExternalDocumentation
	if err := yaml.Unmarshal([]byte(yml), &externalDocumentation); err != nil {
		t.Fatal(err)
	}
	if externalDocumentation.description != "Find more info here" {
		t.Errorf("unexpected externalDocumentation.description: %s", externalDocumentation.description)
		return
	}
	if externalDocumentation.url != "https://example.com" {
		t.Errorf("unexpected externalDocumentation.url: %s", externalDocumentation.url)
		return
	}
}

func TestParameterUnmarshalYAML(t *testing.T) {
	t.Run("header parameter", func(t *testing.T) {
		yml := `name: token
in: header
description: token to be passed as a header
required: true
schema:
  type: array
  items:
    type: integer
    format: int64
style: simple`
		var parameter Parameter
		if err := yaml.Unmarshal([]byte(yml), &parameter); err != nil {
			t.Fatal(err)
		}
		if parameter.name != "token" {
			t.Errorf("unexpected paramater.name: %s", parameter.name)
			return
		}
		if parameter.in != "header" {
			t.Errorf("unexpected parameter.in: %s", parameter.in)
			return
		}
		if parameter.description != "token to be passed as a header" {
			t.Errorf("unexpected parameter.description: %s", parameter.description)
			return
		}
		if parameter.required != true {
			t.Errorf("unexpected parameter.required: %t", parameter.required)
			return
		}
		if parameter.schema.type_ != "array" {
			t.Errorf("unexpected parameter.schema.type: %s", parameter.schema.type_)
			return
		}
		if parameter.schema.items.type_ != "integer" {
			t.Errorf("unexpected parameter.schema.items.type: %s", parameter.schema.items.type_)
			return
		}
		if parameter.schema.items.format != "int64" {
			t.Errorf("unexpected parameter.schema.items.format: %s", parameter.schema.items.format)
			return
		}
		if parameter.style != "simple" {
			t.Errorf("unexpected paarameter.style: %s", parameter.style)
			return
		}
	})
	t.Run("path parameter", func(t *testing.T) {
		yml := `name: username
in: path
description: username to fetch
required: true
schema:
  type: string`
		var parameter Parameter
		if err := yaml.Unmarshal([]byte(yml), &parameter); err != nil {
			t.Fatal(err)
		}
		if parameter.name != "username" {
			t.Errorf("unexpected parameter.name: %s", parameter.name)
			return
		}
		if parameter.in != "path" {
			t.Errorf("unexpected parameter.in: %s", parameter.in)
			return
		}
		if parameter.description != "username to fetch" {
			t.Errorf("unexpected parameter.description: %s", parameter.description)
			return
		}
		if parameter.required != true {
			t.Errorf("unexpected parameter.required: %t", parameter.required)
			return
		}
		if parameter.schema.type_ != "string" {
			t.Errorf("unexpected parameter.schema.type: %s", parameter.schema.type_)
			return
		}
	})
	t.Run("optional query parameter", func(t *testing.T) {
		yml := `name: id
in: query
description: ID of the object to fetch
required: false
schema:
  type: array
  items:
    type: string
style: form
explode: true`
		var parameter Parameter
		if err := yaml.Unmarshal([]byte(yml), &parameter); err != nil {
			t.Fatal(err)
		}
		if parameter.name != "id" {
			t.Errorf("unexpected parameter.name: %s", parameter.name)
			return
		}
		if parameter.in != "query" {
			t.Errorf("unexpected parameter.in: %s", parameter.in)
			return
		}
		if parameter.description != "ID of the object to fetch" {
			t.Errorf("unexpected parameter.description: %s", parameter.description)
			return
		}
		if parameter.required != false {
			t.Errorf("unexpected parameter.required: %t", parameter.required)
			return
		}
		if parameter.schema.type_ != "array" {
			t.Errorf("unexpected parameter.schema.type: %s", parameter.schema.type_)
			return
		}
		if parameter.schema.items.type_ != "string" {
			t.Errorf("unexpected parameter.schema.items.type: %s", parameter.schema.items.type_)
			return
		}
		if parameter.style != "form" {
			t.Errorf("unexpected parameter.style: %s", parameter.style)
			return
		}
		if parameter.explode != true {
			t.Errorf("unexpected parameter.explode: %t", parameter.explode)
			return
		}
	})
	t.Run("free form", func(t *testing.T) {
		yml := `in: query
name: freeForm
schema:
  type: object
  additionalProperties:
    type: integer
style: form`
		var parameter Parameter
		if err := yaml.Unmarshal([]byte(yml), &parameter); err != nil {
			t.Fatal(err)
		}
		if parameter.in != "query" {
			t.Errorf("unexpected parameter.in: %s", parameter.in)
			return
		}
		if parameter.name != "freeForm" {
			t.Errorf("unexpected parameter.name: %s", parameter.name)
			return
		}
		if parameter.schema.type_ != "object" {
			t.Errorf("unexpected parameter.schema.type_: %s", parameter.schema.type_)
			return
		}
		if parameter.schema.additionalProperties.type_ != "integer" {
			t.Errorf("unexpected parameter.schema.additionalProperties.type_: %s", parameter.schema.additionalProperties.type_)
			return
		}
		if parameter.style != "form" {
			t.Errorf("unexpected parameter.style: %s", parameter.style)
			return
		}
	})
	t.Run("complex parameter", func(t *testing.T) {
		yml := `in: query
name: coordinates
content:
  application/json:
    schema:
      type: object
      required:
        - lat
        - long
      properties:
        lat:
          type: number
        long:
          type: number`
		var parameter Parameter
		if err := yaml.Unmarshal([]byte(yml), &parameter); err != nil {
			t.Fatal(err)
		}
		if parameter.in != "query" {
			t.Errorf("unexpected parameter.in: %s", parameter.in)
			return
		}
		if parameter.name != "coordinates" {
			t.Errorf("unexpected parameter.name: %s", parameter.name)
			return
		}
		if parameter.content["application/json"].schema.type_ != "object" {
			t.Errorf("unexpected parameter.content.application/json.schema.type_: %s", parameter.content["application/json"].schema.type_)
			return
		}
		if parameter.content["application/json"].schema.required[0] != "lat" {
			t.Errorf("unexpected parameter.content.application/json.schema.required.0: %s", parameter.content["application/json"].schema.required[0])
			return
		}
		if parameter.content["application/json"].schema.required[1] != "long" {
			t.Errorf("unexpected parameter.content.application/json..schema.required.1: %s", parameter.content["application/json"].schema.required[1])
			return
		}
		if parameter.content["application/json"].schema.properties["lat"].type_ != "number" {
			t.Errorf("unexpected parameter.content.application/json.schema.properties.lat.type: %s", parameter.content["application/json"].schema.properties["lat"].type_)
			return
		}
		if parameter.content["application/json"].schema.properties["long"].type_ != "number" {
			t.Errorf("unexpected parameter.content.application/json.schema.properties.long.type: %s", parameter.content["application/json"].schema.properties["long"].type_)
			return
		}
	})
}

func TestRequestBodyUnmarshalYAML(t *testing.T) {
	t.Run("with a referenced model", func(t *testing.T) {
		yml := `description: user to add to the system
content:
  'application/json':
    schema:
      $ref: '#/components/schemas/User'
    examples:
      user:
        summary: User Example
        externalValue: 'http://foo.bar/examples/user-example.json'
  'application/xml':
    schema:
      $ref: '#/components/schemas/User'
    examples:
      user:
        summary: User Example in XML
        externalValue: 'http://foo.bar/examples/user-example.xml'
  'text/plain':
    examples:
      user:
        summary: User example in text plain format
        externalValue: 'http://foo.bar/examples/user-example.txt'
  '*/*':
    examples:
      user:
        summary: User example in other format
        externalValue: 'http://foo.bar/examples/user-example.whatever'`
		var requestBody RequestBody
		if err := yaml.Unmarshal([]byte(yml), &requestBody); err != nil {
			t.Fatal(err)
		}
		if requestBody.description != "user to add to the system" {
			t.Errorf("unexpected requestBody.description: %s", requestBody.description)
			return
		}
		t.Run("application/json", func(t *testing.T) {
			mediaType, ok := requestBody.content["application/json"]
			if !ok {
				t.Error("requestBody.content.application/json is not found")
				return
			}
			if mediaType.schema.reference != "#/components/schemas/User" {
				t.Errorf("unexpected requestBody.content.application/json.schema.$ref: %s", mediaType.schema.reference)
				return
			}
			example, ok := mediaType.examples["user"]
			if !ok {
				t.Error("requestBody.content.application/json.examples.user is not found")
				return
			}
			if example.summary != "User Example" {
				t.Errorf("unexpected requestBody.content.application/json.examples.user.summary: %s", example.summary)
				return
			}
			if example.externalValue != "http://foo.bar/examples/user-example.json" {
				t.Errorf("unexpected requestBody.content.application/json.examples.user.externalValue: %s", example.externalValue)
				return
			}
		})
		t.Run("application/xml", func(t *testing.T) {
			mediaType, ok := requestBody.content["application/xml"]
			if !ok {
				t.Error("requestBody.content.application/xml is not found")
				return
			}
			if mediaType.schema.reference != "#/components/schemas/User" {
				t.Errorf("unexpected requestBody.content.application/xml.schema.$ref: %s", mediaType.schema.reference)
				return
			}
			example, ok := mediaType.examples["user"]
			if !ok {
				t.Error("requestBody.content.application/xml.examples.user is not found")
				return
			}
			if example.summary != "User Example in XML" {
				t.Errorf("unexpected requestBody.content.application/xml.examples.user.summary: %s", example.summary)
				return
			}
			if example.externalValue != "http://foo.bar/examples/user-example.xml" {
				t.Errorf("unexpected requestBody.content.application/xml.examples.user.externalValue: %s", example.externalValue)
				return
			}
		})
		t.Run("text/plain", func(t *testing.T) {
			mediaType, ok := requestBody.content["text/plain"]
			if !ok {
				t.Error("requestBody.content.text/plain is not found")
				return
			}
			example, ok := mediaType.examples["user"]
			if !ok {
				t.Error("requestBody.content.text/plain.examples.user is not found")
				return
			}
			if example.summary != "User example in text plain format" {
				t.Errorf("unexpected requestBody.content.text/plain.examples.user.summary: %s", example.summary)
				return
			}
			if example.externalValue != "http://foo.bar/examples/user-example.txt" {
				t.Errorf("unexpected requestBody.content.text/plain.examples.user.externalValue: %s", example.externalValue)
				return
			}
		})
		t.Run("*/*", func(t *testing.T) {
			mediaType, ok := requestBody.content["*/*"]
			if !ok {
				t.Error("requestBody.content.*/* is not found")
				return
			}
			example, ok := mediaType.examples["user"]
			if !ok {
				t.Error("requestBody.content.*/*.examples.user is not found")
				return
			}
			if example.summary != "User example in other format" {
				t.Errorf("unexpected requestBody.content.*/*.examples.user.summary: %s", example.summary)
				return
			}
			if example.externalValue != "http://foo.bar/examples/user-example.whatever" {
				t.Errorf("unexpected requestBody.content.*/*.examples.user.externalValue: %s", example.externalValue)
				return
			}
		})
	})
	t.Run("array of string", func(t *testing.T) {
		yml := `description: user to add to the system
required: true
content:
  text/plain:
    schema:
      type: array
      items:
        type: string`
		var requestBody RequestBody
		if err := yaml.Unmarshal([]byte(yml), &requestBody); err != nil {
			t.Fatal(err)
		}
		if requestBody.description != "user to add to the system" {
			t.Errorf("unexpected requestBody.description: %s", requestBody.description)
			return
		}
		if requestBody.required != true {
			t.Errorf("unexpected requestBody.required: %t", requestBody.required)
			return
		}
		mediaType, ok := requestBody.content["text/plain"]
		if !ok {
			t.Error("requestBody.content.text/plain is not found")
			return
		}
		if mediaType.schema.type_ != "array" {
			t.Errorf("unexpected mediaType.schema.type: %s", mediaType.schema.type_)
			return
		}
		if mediaType.schema.items.type_ != "string" {
			t.Errorf("unexpected mediaType.schema.items.type: %s", mediaType.schema.items.type_)
			return
		}
	})
}

func TestMediaTypeUnmarshalYAML(t *testing.T) {
	yml := `application/json:
  schema:
    $ref: "#/components/schemas/Pet"
  examples:
    cat:
      summary: An example of a cat
      value:
        name: Fluffy
        petType: Cat
        color: White
        gender: male
        breed: Persian
    dog:
      summary: An example of a dog with a cat's name
      value:
        name: Puma
        petType: Dog
        color: Black
        gender: Female
        breed: Mixed
    frog:
      $ref: "#/components/examples/frog-example"`
	var target map[string]*MediaType
	if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
		t.Fatal(err)
	}
	mediaType, ok := target["application/json"]
	if !ok {
		t.Error("application/json is not found")
		return
	}
	if mediaType.schema.reference != "#/components/schemas/Pet" {
		t.Errorf("unexpected mediaType.schema.$ref: %s", mediaType.schema.reference)
		return
	}
	t.Run("cat", func(t *testing.T) {
		example, ok := mediaType.examples["cat"]
		if !ok {
			t.Error("mediaType.examples.cat is not found")
			return
		}
		value, ok := example.value.(map[string]interface{})
		if !ok {
			t.Errorf("mediaType.examples.cat.value is assumed map[string]interface but %v", reflect.TypeOf(example.value))
			return
		}
		if name, ok := value["name"]; !ok {
			t.Error("mediaType.examples.cat.value.name is not found")
			return
		} else if name != "Fluffy" {
			t.Errorf("unexpected mediaType.examples.cat.value.name: %s", name)
		}
		if petType, ok := value["petType"]; !ok {
			t.Error("mediaType.examples.cat.value.petType is not found")
			return
		} else if petType != "Cat" {
			t.Errorf("unexpected mediaType.examples.cat.value.petType: %s", petType)
		}
		if color, ok := value["color"]; !ok {
			t.Error("mediaType.examples.cat.value.color is not found")
			return
		} else if color != "White" {
			t.Errorf("unexpected mediaType.examples.cat.value.color: %s", color)
		}
		if gender, ok := value["gender"]; !ok {
			t.Error("mediaType.examples.cat.value.gender is not found")
			return
		} else if gender != "male" {
			t.Errorf("unexpected mediaType.examples.cat.value.gender: %s", gender)
		}
		if breed, ok := value["breed"]; !ok {
			t.Error("mediaType.examples.cat.value.breed is not found")
			return
		} else if breed != "Persian" {
			t.Errorf("unexpected mediaType.examples.cat.value.breed: %s", breed)
			return
		}
	})
	t.Run("dog", func(t *testing.T) {
		example, ok := mediaType.examples["dog"]
		if !ok {
			t.Error("mediaType.examples.dog is not found")
			return
		}
		value, ok := example.value.(map[string]interface{})
		if !ok {
			t.Errorf("mediaType.examples.dog.value is assumed map[string]interface but %v", reflect.TypeOf(example.value))
			return
		}
		if name, ok := value["name"]; !ok {
			t.Error("mediaType.examples.dog.value.name is not found")
			return
		} else if name != "Puma" {
			t.Errorf("unexpected mediaType.examples.dog.value.name: %s", name)
		}
		if petType, ok := value["petType"]; !ok {
			t.Error("mediaType.examples.dog.value.petType is not found")
			return
		} else if petType != "Dog" {
			t.Errorf("unexpected mediaType.examples.dog.value.petType: %s", petType)
		}
		if color, ok := value["color"]; !ok {
			t.Error("mediaType.examples.dog.value.color is not found")
			return
		} else if color != "Black" {
			t.Errorf("unexpected mediaType.examples.dog.value.color: %s", color)
		}
		if gender, ok := value["gender"]; !ok {
			t.Error("mediaType.examples.dog.value.gender is not found")
			return
		} else if gender != "Female" {
			t.Errorf("unexpected mediaType.examples.dog.value.gender: %s", gender)
		}
		if breed, ok := value["breed"]; !ok {
			t.Error("mediaType.examples.dog.value.breed is not found")
			return
		} else if breed != "Mixed" {
			t.Errorf("unexpected mediaType.examples.dog.value.breed: %s", breed)
			return
		}
	})
	t.Run("frog", func(t *testing.T) {
		example, ok := mediaType.examples["frog"]
		if !ok {
			t.Error("mediaType.examples.frog is not found")
			return
		}
		if example.reference != "#/components/examples/frog-example" {
			t.Errorf("unexpected mediaType.examples.frog.$ref: %s", example.reference)
			return
		}
	})
}

func TestEncodingUnmarshalYAML(t *testing.T) {
	yml := `requestBody:
  content:
    multipart/mixed:
      schema:
        type: object
        properties:
          id:
            # default is text/plain
            type: string
            format: uuid
          address:
            # default is application/json
            type: object
            properties: {}
          historyMetadata:
            # need to declare XML format!
            description: metadata in XML format
            type: object
            properties: {}
          profileImage:
            # default is application/octet-stream, need to declare an image type only!
            type: string
            format: binary
      encoding:
        historyMetadata:
          # require XML Content-Type in utf-8 encoding
          contentType: application/xml; charset=utf-8
        profileImage:
          # only accept png/jpeg
          contentType: image/png, image/jpeg
          headers:
            X-Rate-Limit-Limit:
              description: The number of allowed requests in the current period
              schema:
                type: integer`
	var target struct {
		RequestBody RequestBody `yaml:"requestBody"`
	}
	if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
		t.Fatal(err)
	}
	if _, ok := target.RequestBody.content["multipart/mixed"]; !ok {
		t.Error("requestBody.content.multipart/mixed is not found")
		return
	}
	schema := target.RequestBody.content["multipart/mixed"].schema
	if schema.type_ != "object" {
		t.Errorf("unexpected requestBody.content.multipart/mixed.schema.type: %s", schema.type_)
		return
	}
	id, ok := schema.properties["id"]
	if !ok {
		t.Error("requestBody.content.multipart/mixed.schema.properties.id is not found")
		return
	}
	if id.type_ != "string" {
		t.Errorf("unexpected id.type: %s", id.type_)
		return
	}
	if id.format != "uuid" {
		t.Errorf("unexpected id.format: %s", id.format)
		return
	}
	address, ok := schema.properties["address"]
	if !ok {
		t.Error("requestBody.content.multipart/mixed.schema.properties.address is not found")
		return
	}
	if address.type_ != "object" {
		t.Errorf("unexpected id.type: %s", id.type_)
		return
	}
	historyMetadata, ok := schema.properties["historyMetadata"]
	if !ok {
		t.Error("requestBody.content.multipart/mixed.schema.properties.historyMetadata is not found")
		return
	}
	if historyMetadata.description != "metadata in XML format" {
		t.Errorf("unexpected historyMetadata.description: %s", historyMetadata.description)
		return
	}
	if historyMetadata.type_ != "object" {
		t.Errorf("unexpected historyMetadata.type: %s", historyMetadata.type_)
		return
	}
	profileImage, ok := schema.properties["profileImage"]
	if !ok {
		t.Error("requestBody.content.multipart/mixed.schema.properties.profileImage is not found")
		return
	}
	if profileImage.type_ != "string" {
		t.Errorf("unexpected id.type: %s", id.type_)
		return
	}
	if profileImage.format != "binary" {
		t.Errorf("unexpected id.format: %s", id.format)
		return
	}
	t.Run("historyMetadata", func(t *testing.T) {
		encoding := target.RequestBody.content["multipart/mixed"].encoding["historyMetadata"]
		if encoding.contentType != "application/xml; charset=utf-8" {
			t.Errorf("unexpected encoding.contentType: %s", encoding.contentType)
			return
		}
	})
	t.Run("profileImage", func(t *testing.T) {
		encoding := target.RequestBody.content["multipart/mixed"].encoding["profileImage"]
		if encoding.contentType != "image/png, image/jpeg" {
			t.Errorf("unexpected encoding.contentType: %s", encoding.contentType)
			return
		}
		xRateLimitLimit, ok := encoding.headers["X-Rate-Limit-Limit"]
		if !ok {
			t.Error("encoding.headers.X-Rate-Limit-Limit is not found")
			return
		}
		if xRateLimitLimit.description != "The number of allowed requests in the current period" {
			t.Errorf("unexpeceted encoding.headers.X-Rate-Limit-Limit.description: %s", xRateLimitLimit.description)
			return
		}
		if xRateLimitLimit.schema.type_ != "integer" {
			t.Errorf("unexpected encoding.headers.X-Rate-Limit-Limit.schema.type: %s", xRateLimitLimit.schema.type_)
			return
		}
	})
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
