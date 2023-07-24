package hqgourl

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct { //nolint:govet // To be refactored.
		name   string
		input  string
		output *URL
		err    error
	}{
		{
			name:  "Test example URL",
			input: "https://sub.example.com:8080/path/to/file.txt",
			output: &URL{
				Domain:      "sub.example.com",
				ETLDPlusOne: "example.com",
				Subdomain:   "sub",
				RootDomain:  "example",
				TLD:         "com",
				Port:        "8080",
				Extension:   ".txt",
			},
			err: nil,
		},
	}

	for index := range tests {
		tt := tests[index]

		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse(%q) returned error %v", tt.input, err)
			}

			if got.Domain != tt.output.Domain ||
				got.ETLDPlusOne != tt.output.ETLDPlusOne ||
				got.Subdomain != tt.output.Subdomain ||
				got.RootDomain != tt.output.RootDomain ||
				got.TLD != tt.output.TLD ||
				got.Port != tt.output.Port ||
				got.Extension != tt.output.Extension {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, got, tt.output)
			}
		})
	}
}

func TestAddDefaultScheme(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		scheme string
		output string
	}{
		{
			name:   "Case: localhost",
			url:    "localhost",
			scheme: "http",
			output: "http://localhost",
		},
		{
			name:   "Case: example.com",
			url:    "example.com",
			scheme: "http",
			output: "http://example.com",
		},
		{
			name:   "Case: //example.com",
			url:    "//example.com",
			scheme: "http",
			output: "http://example.com",
		},
		{
			name:   "Case: ://example.com",
			url:    "://example.com",
			scheme: "http",
			output: "http://example.com",
		},
		{
			name:   "Case: https://example.com",
			url:    "https://example.com",
			scheme: "http",
			output: "https://example.com",
		},
	}

	for index := range tests {
		tt := tests[index]

		t.Run(tt.name, func(t *testing.T) {
			got := AddDefaultScheme(tt.url, tt.scheme)

			if got != tt.output {
				t.Errorf("AddDefaultScheme(%q, %q) = %v, want %v", tt.url, tt.scheme, got, tt.output)
			}
		})
	}
}

func TestSplitHost(t *testing.T) {
	tests := []struct {
		name   string
		host   string
		domain string
		port   string
	}{
		{
			name:   "Case: localhost",
			host:   "localhost",
			domain: "localhost",
			port:   "",
		},
		{
			name:   "Case: example.com",
			host:   "example.com",
			domain: "example.com",
			port:   "",
		},
		{
			name:   "Case: localhost:8080",
			host:   "localhost:8080",
			domain: "localhost",
			port:   "8080",
		},
		{
			name:   "Case: example.com:8080",
			host:   "example.com:8080",
			domain: "example.com",
			port:   "8080",
		},
	}

	for index := range tests {
		tt := tests[index]

		t.Run(tt.name, func(t *testing.T) {
			domain, port := SplitHost(tt.host)

			if domain != tt.domain || port != tt.port {
				t.Errorf("SplitHost(%q) = %v, %v, want %v, %v", tt.host, domain, port, tt.domain, tt.port)
			}
		})
	}
}
