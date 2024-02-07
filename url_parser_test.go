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
			"http://example.com",
			"http",
			&hqgourl.URL{
				URL: &url.URL{
					Scheme: "http",
					Host:   "example.com",
				},
				Domain: &hqgourl.Domain{
					Sub:      "",
					Root:     "example",
					TopLevel: "com",
				},
			},
			false,
		},
		{
			"example.com",
			"http",
			&hqgourl.URL{
				URL: &url.URL{
					Scheme: "http",
					Host:   "example.com",
				},
				Domain: &hqgourl.Domain{
					Sub:      "",
					Root:     "example",
					TopLevel: "com",
				},
			},
			false,
		},
		{
			"http://example.com/path/file.html",
			"http",
			&hqgourl.URL{
				URL: &url.URL{
					Scheme: "http",
					Host:   "example.com",
					Path:   "/path/file.html",
				},
				Domain: &hqgourl.Domain{
					Sub:      "",
					Root:     "example",
					TopLevel: "com",
				},
				Extension: ".html",
			},
			false,
		},
	}

	for _, c := range cases {
		c := c

		t.Run(fmt.Sprintf("Parse(%q)", c.rawURL), func(t *testing.T) {
			t.Parallel()

			parser := hqgourl.NewURLParser(
				hqgourl.URLParserWithDefaultScheme(c.defaultScheme),
			)

			parsedURL, err := parser.Parse(c.rawURL)

			if (err != nil) != c.expectParseErr {
				t.Errorf("Parse(%q) error = %v, expectParseErr %v", c.rawURL, err, c.expectParseErr)

				return
			}

			if !reflect.DeepEqual(parsedURL, c.expectedParsedURL) {
				t.Errorf("Parse(%q) = %+v, want %+v", c.rawURL, parsedURL, c.expectedParsedURL)
			}
		})
	}
}
