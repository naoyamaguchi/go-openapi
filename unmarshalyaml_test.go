package openapi

import (
	"os"
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
			Servers []*Server `yaml:"servers"`
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

func mustOpen(filename string) *os.File {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return f
}
