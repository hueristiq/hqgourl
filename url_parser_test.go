package hqgourl_test

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/hueristiq/hqgourl"
)

func TestNewURLParser(t *testing.T) {
	t.Parallel()

	parser := hqgourl.NewURLParser()

	if parser == nil {
		t.Error("NewURLParser returned nil")
	}

	scheme := "https"

	parserWithDefaultScheme := hqgourl.NewURLParser(hqgourl.URLParserWithDefaultScheme(scheme))

	expectedDefaultScheme := parserWithDefaultScheme.DefaultScheme()

	if expectedDefaultScheme != scheme {
		t.Errorf("Expected default scheme '%s', got '%s'", scheme, expectedDefaultScheme)
	}
}

func TestURLParser_Parse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		rawURL         string
		defaultScheme  string
		expectedURL    *hqgourl.URL
		expectParseErr bool
	}{
		{
			name:          "Basic URL",
			rawURL:        "http://example.com",
			defaultScheme: "http",
			expectedURL: &hqgourl.URL{
				URL: &url.URL{
					Scheme: "http",
					Host:   "example.com",
				},
				RootDomain:     "example",
				TopLevelDomain: "com",
			},
		},
		{
			name:          "Basic URL without scheme",
			rawURL:        "example.com",
			defaultScheme: "http",
			expectedURL: &hqgourl.URL{
				URL: &url.URL{
					Scheme: "http",
					Host:   "example.com",
				},
				RootDomain:     "example",
				TopLevelDomain: "com",
			},
		},
		{
			name:          "Standard URL with http",
			rawURL:        "http://example.com/path/file.html",
			defaultScheme: "http",
			expectedURL: &hqgourl.URL{
				URL: &url.URL{
					Scheme: "http",
					Host:   "example.com",
					Path:   "/path/file.html",
				},
				RootDomain:     "example",
				TopLevelDomain: "com",
				Extension:      ".html",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			parser := hqgourl.NewURLParser(
				hqgourl.URLParserWithDefaultScheme(tc.defaultScheme),
			)

			parsedURL, err := parser.Parse(tc.rawURL)

			if (err != nil) != tc.expectParseErr {
				t.Errorf("Parse() error = %v, expectParseErr %v", err, tc.expectParseErr)

				return
			}

			if !reflect.DeepEqual(parsedURL, tc.expectedURL) {
				t.Errorf("Parse() = %+v, want %+v", parsedURL, tc.expectedURL)
			}
		})
	}
}
