package openapi

import (
	"regexp"
	"strings"
)

var (
	urlTemplateVarRegexp = regexp.MustCompile("{[^}]+}") // nolint[gocheckonglobals]

	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$") //nolint[lll]
)

type raw []byte

func (v *raw) UnmarshalYAML(b []byte) error {
	*v = b
	return nil
}

func extension(proxy map[string]raw) map[string]interface{} {
	extension := map[string]interface{}{}
	for k, v := range proxy {
		if strings.HasPrefix(k, "x-") {
			extension[k] = v
		}
	}
	if len(extension) == 0 {
		return nil
	}
	return extension
}

func isOneOf(s string, list []string) bool {
	for _, t := range list {
		if t == s {
			return true
		}
	}
	return false
}

const (
	rfc5234Alpha     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rfc5234Digit     = "0123456789"
	rfc7230TChar     = "!#$%&'*+-.^_`|~" + rfc5234Digit + rfc5234Alpha
	rfc7159Unescaped = "\x20\x21\x23\x24\x25\x26\x27\x28\x29\x2a\x2b\x2c\x2d\x2e\x2f\x30"
)

func matchRuntimeExpr(expr string) bool {
	if !strings.HasPrefix(expr, "$") {
		return false
	}
	if expr == "$url" || expr == "$method" || expr == "$statusCode" {
		return true
	}
	var source string
	if !strings.HasPrefix(expr, "$request.") {
		if !strings.HasPrefix(expr, "$response.") {
			return false
		}
		source = strings.TrimPrefix(expr, "$response.")
	} else {
		source = strings.TrimPrefix(expr, "$request.")
	}
	if len(source) == 0 {
		return false
	}
	var name string
	switch {
	case strings.HasPrefix(source, "header."):
		token := strings.TrimPrefix(source, "header.")
		if len(token) == 0 {
			return false
		}
		return len(strings.Trim(token, rfc7230TChar)) == 0
	case strings.HasPrefix(source, "body"):
		if strings.Contains(source, "#") {
			split := strings.Split(source, "#")
			if split[0] != "body" {
				return false
			}
			fragment := split[1]
			return strings.HasPrefix(fragment, "/")
		}
	case strings.HasPrefix(source, "query."):
		name = strings.TrimPrefix(source, "query.")
	case strings.HasPrefix(source, "path."):
		name = strings.TrimPrefix(source, "path.")
	}
	return len(name) != 0
}
