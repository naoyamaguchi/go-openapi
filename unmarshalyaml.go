package openapi

import "strings"

type raw []byte

func (v *raw) UnmarshalYAML(b []byte) error {
	*v = b
	return nil
}

func extension(proxy map[string]raw) Extension {
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
