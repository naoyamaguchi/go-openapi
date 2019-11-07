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
