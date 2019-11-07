package openapi

import (
	"strconv"
	"testing"
)

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
