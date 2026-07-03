package main

import (
	"net/http"
	"testing"
)

func TestAPIKeyFromRequest(t *testing.T) {
	cases := []struct {
		name    string
		headers map[string]string
		want    string
	}{
		{name: "no headers", headers: nil, want: ""},
		{name: "dedicated header wins", headers: map[string]string{
			"X-Sentiary-User-Api-Key": "dedicated",
			"Authorization":           "Bearer other",
		}, want: "dedicated"},
		{name: "dedicated header trimmed", headers: map[string]string{
			"X-Sentiary-User-Api-Key": "  spaced  ",
		}, want: "spaced"},
		{name: "bearer scheme", headers: map[string]string{"Authorization": "Bearer abc123"}, want: "abc123"},
		{name: "bearer lowercase scheme", headers: map[string]string{"Authorization": "bearer abc123"}, want: "abc123"},
		{name: "ribbon scheme", headers: map[string]string{"Authorization": "Ribbon abc123"}, want: "abc123"},
		{name: "bearer extra spaces", headers: map[string]string{"Authorization": "Bearer   abc123"}, want: "abc123"},
		{name: "bare token", headers: map[string]string{"Authorization": "rawtoken"}, want: "rawtoken"},
		{name: "lone bearer keyword falls back", headers: map[string]string{"Authorization": "Bearer"}, want: ""},
		{name: "lone ribbon keyword falls back", headers: map[string]string{"Authorization": "ribbon"}, want: ""},
		{name: "bearer with only spaces falls back", headers: map[string]string{"Authorization": "Bearer   "}, want: ""},
		{name: "empty authorization", headers: map[string]string{"Authorization": "   "}, want: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := &http.Request{Header: http.Header{}}
			for key, value := range tc.headers {
				request.Header.Set(key, value)
			}
			if got := apiKeyFromRequest(request); got != tc.want {
				t.Fatalf("apiKeyFromRequest() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestAPIKeyFromRequestNil(t *testing.T) {
	if got := apiKeyFromRequest(nil); got != "" {
		t.Fatalf("apiKeyFromRequest(nil) = %q, want empty", got)
	}
	if got := apiKeyFromRequest(&http.Request{}); got != "" {
		t.Fatalf("apiKeyFromRequest(no header) = %q, want empty", got)
	}
}
