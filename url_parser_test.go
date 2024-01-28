package hqgourl_test

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/hueristiq/hqgourl"
)

func TestNewURLParser(t *testing.T) {
	t.Parallel()

	parser := hqgourl.NewURLParser()

	if parser == nil {
		t.Error("NewURLParser() = nil; want non-nil")
	}

	scheme := "https"

	parserWithDefaultScheme := hqgourl.NewURLParser(hqgourl.URLParserWithDefaultScheme(scheme))

	if parser == nil {
		t.Errorf("NewURLParser(URLParserWithDefaultScheme(%s)) = nil; want non-nil", scheme)
	}

	expectedDefaultScheme := parserWithDefaultScheme.DefaultScheme()

	if expectedDefaultScheme != scheme {
		t.Errorf("NewURLParser(URLParserWithDefaultScheme(%s)).DefaultScheme() = '%s', want '%s'", scheme, expectedDefaultScheme, scheme)
	}
}

func TestURLParser_Parse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		rawURL            string
		defaultScheme     string
		expectedParsedURL *hqgourl.URL
		expectParseErr    bool
	}{
		{
			rawURL:        "http://example.com",
			defaultScheme: "http",
			expectedParsedURL: &hqgourl.URL{
				URL: &url.URL{
					Scheme: "http",
					Host:   "example.com",
				},
				RootDomain:     "example",
				TopLevelDomain: "com",
			},
		},
		{
			rawURL:        "example.com",
			defaultScheme: "http",
			expectedParsedURL: &hqgourl.URL{
				URL: &url.URL{
					Scheme: "http",
					Host:   "example.com",
				},
				RootDomain:     "example",
				TopLevelDomain: "com",
			},
		},
		{
			rawURL:        "http://example.com/path/file.html",
			defaultScheme: "http",
			expectedParsedURL: &hqgourl.URL{
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

	for _, c := range cases {
		c := c

		t.Run(fmt.Sprintf("Parse(%s)", c.rawURL), func(t *testing.T) {
			t.Parallel()

			parser := hqgourl.NewURLParser(
				hqgourl.URLParserWithDefaultScheme(c.defaultScheme),
			)

			parsedURL, err := parser.Parse(c.rawURL)

			if (err != nil) != c.expectParseErr {
				t.Errorf("Parse() error = %v, expectParseErr %v", err, c.expectParseErr)

				return
			}

			if !reflect.DeepEqual(parsedURL, c.expectedParsedURL) {
				t.Errorf("Parse() = %+v, want %+v", parsedURL, c.expectedParsedURL)
			}
		})
	}
}
