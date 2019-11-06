//go:generate go run mkunmarshalyaml.go astlib.go
//go:generate go run mkgetter.go astlib.go

package openapi

//+object
// OpenAPI is the root document object of the OpenAPI document.
type OpenAPI struct {
	openapi      string `required:"yes"`
	info         *Info  `required:"yes"`
	servers      []*Server
	paths        *Paths `required:"yes"`
	components   *Components
	security     []*SecurityRequirement
	tags         []*Tag
	externalDocs *ExternalDocumentation

	extension Extension
}

//+object
// Info provides metadata about the API.
type Info struct {
	title          string `required:"yes"`
	description    string
	termsOfService string
	contact        *Contact
	license        *License
	version        string `required:"yes"`

	extension Extension
}

//+object
// Contact information for the exposed API.
type Contact struct {
	name  string
	url   string
	email string

	extension Extension
}

//+object
// License information for the exposed API.
type License struct {
	name string `required:"yes"`
	url  string

	extension Extension
}

//+object
// Server is an object representing a Server.
type Server struct {
	url         string `required:"yes"`
	description string
	variables   map[string]*ServerVariable

	extension Extension
}

//+object
// ServerVariable is an object representing a Server Variable for serverURL template substitution.
type ServerVariable struct {
	enum []string
	//nolint[golint]
	default_    string `required:"yes" yaml:"default"`
	description string

	extension Extension
}

//+object
// Components holds a set of reusable objects for different aspects of the OAS.
type Components struct {
	schemas         map[string]*Schema
	responses       map[string]*Response
	parameters      map[string]*Parameter
	examples        map[string]*Example
	requestBodies   map[string]*RequestBody
	headers         map[string]*Header
	securitySchemes map[string]*SecurityScheme
	links           map[string]*Link
	callbacks       map[string]*Callback

	extension Extension
}

//+object
// Paths holds the relative paths to the individual endpoints and their operations.
type Paths struct {
	paths map[string]*PathItem `yaml:",inline"`

	extension Extension
}

//+object
// PathItem describes the operations available on a single path.
type PathItem struct {
	summary     string
	description string
	get         *Operation
	put         *Operation
	post        *Operation
	delete      *Operation
	options     *Operation
	head        *Operation
	patch       *Operation
	trace       *Operation
	servers     []*Server
	parameters  []*Parameter

	extension Extension
}

//+object
// Operation describes a single API operation on a path.
type Operation struct {
	tags         []string
	summary      string
	description  string
	externalDocs *ExternalDocumentation
	operationID  string
	parameters   *Parameter
	requestBody  *RequestBody
	responses    *Responses `required:"yes"`
	callbacks    map[string]*Callback
	deprecated   bool
	security     []*SecurityRequirement
	servers      []*Server

	extension Extension
}

//+object
// ExternalDocumentation allows referencing an external resource for extended documentation.
type ExternalDocumentation struct {
	description string
	url         string `required:"yes"`

	extension Extension
}

//+object
// Parameter describes a single operation parameter.
type Parameter struct {
	name            string `required:"yes"`
	in              string `required:"yes"`
	description     string
	required        bool
	deprecated      bool
	allowEmptyValue bool
	style           string
	explode         bool
	allowReserved   bool
	schema          *Schema
	example         interface{}
	examples        map[string]*Example
	content         map[string]*MediaType

	extension Extension
}

//+object
// RequestBody describes a single request body.
type RequestBody struct {
	description string
	content     map[string]*MediaType `required:"yes"`
	required    bool

	extension Extension
}

//+object
// MediaType provides schema and examples for the media type identified by its key.
type MediaType struct {
	schema   *Schema
	example  interface{}
	examples map[string]*Example
	encoding map[string]*Encoding

	extension Extension
}

//+object
// Encoding is a single encoding definition applied to a single schema property.
type Encoding struct {
	contentType   string
	headers       map[string]*Header
	style         string
	explode       string
	allowReserved bool

	extension Extension
}

//+object
// Responses is a container for the expected responses of an operation.
type Responses struct {
	responses map[string]*Response

	extension Extension
}

//+object
// Response describes a single response from an API Operation, including design-time,
//+object
// static links to operations based on the response.
type Response struct {
	description string `required:"yes"`
	headers     map[string]*Header
	content     map[string]*MediaType
	links       map[string]*Link

	extension Extension
}

//+object
// Callback is a map of possible out-of band callbacks relatedd to the parent operation.
type Callback struct {
	callback map[string]*PathItem

	extension Extension
}

//+object
// Example object represents an example.
type Example struct {
	summary      string
	description  string
	value        interface{}
	externalVale string

	extension Extension
}

//+object
// Link represents a possible design-time link for a response.
type Link struct {
	operationRef string
	operationId  string
	parameters   map[string]interface{}
	requestBody  interface{}
	description  string
	server       *Server

	extension Extension
}

//+object
// Header object
type Header struct {
	name            string
	in              string
	description     string
	required        bool
	deprecated      bool
	allowEmptyValue bool
	style           string
	explode         bool
	allowReserved   bool
	schema          *Schema
	example         interface{}
	examples        map[string]*Example
	content         map[string]*MediaType

	extension Extension
}

//+object
// Tag adds metadata to a single tag that is used by the Operation Object.
type Tag struct {
	name         string `required:"yes"`
	description  string
	externalDocs *ExternalDocumentation

	extension Extension
}

//+object
// Schema allows the definition of input and output data types.
type Schema struct {
	title            string
	multipleOf       int
	maximum          int
	exclusiveMaximum bool
	minimum          int
	exclusiveMinimum bool
	maxLength        int
	minLength        int
	pattern          string
	maxItems         int
	minItems         int
	maxProperties    int
	minProperties    int
	required         []string
	enum             []string

	type_                string `yaml:"type"`
	allOf                []*Schema
	oneOf                []*Schema
	anyOf                []*Schema
	not                  *Schema
	items                *Schema
	properties           map[string]*Schema
	additionalProperties *Schema
	description          string
	format               string
	default_             string `yaml:"default"`

	nullable      bool
	discriminator *Discriminator
	readOnly      bool
	writeOnly     bool
	xml           *XML
	externalDocs  *ExternalDocumentation
	example       interface{}
	deprecated    bool

	extension Extension
}

//+object
// Discriminator object.
type Discriminator struct {
	propertyName string
	mapping      map[string]string
}

//+object
// XML is a metadata object that allows for more fine-tuned XML model definitions.
type XML struct {
	name      string
	namespace string
	prefix    string
	attribute bool
	wrapped   bool

	extension Extension
}

//+object
// SecuritySchema defines a security scheme that can be used by the operations.
type SecurityScheme struct {
	type_            string `yaml:"type"`
	description      string
	name             string
	in               string
	scheme           string
	bearerFormat     string
	flows            *OAuthFlows
	openIDConnectURL string `yaml:"openIdConnectUrl"`

	extension Extension
}

//+object
// OAuthFlows allows configuration of the supported OAuthFlows.
type OAuthFlows struct {
	implicit          *OAuthFlow
	password          *OAuthFlow
	clientCredentials *OAuthFlow
	authorizationCode *OAuthFlow

	extension Extension
}

//+object
// OAuthFlow is configuration details for a supported OAuth Flow.
type OAuthFlow struct {
	authorizationURL string `yaml:"authorizationUrl"`
	tokenURL         string `yaml:"tokenUrl"`
	refreshURL       string `yaml:"refreshUrl"`
	scopes           map[string]string

	extension Extension
}

//+object
// SecurityRequirements is lists the required security schemes to execute this operation.
type SecurityRequirement map[string][]string

type Extension map[string]interface{}
