package hqgourl_test

import (
	"fmt"
	"testing"

	"github.com/hueristiq/hqgourl"
)

func TestNewDomainParser(t *testing.T) {
	t.Parallel()

	parser := hqgourl.NewDomainParser()
	if parser == nil {
		t.Error("NewDomainParser() = nil; want non-nil")
	}
}

func TestDomainParsing(t *testing.T) {
	t.Parallel()

	cases := []struct {
		domain                 string
		expectedSubdomain      string
		expectedRootDomain     string
		expectedTopLevelDomain string
	}{
		{"www.example.com", "www", "example", "com"},
		{"subdomain.example.co.uk", "subdomain", "example", "co.uk"},
		{"example.com", "", "example", "com"},
		{"localhost", "", "localhost", ""},
	}

	parser := hqgourl.NewDomainParser()

	for _, c := range cases {
		c := c

		t.Run(fmt.Sprintf("Parse(%s)", c.domain), func(t *testing.T) {
			t.Parallel()

			sub, root, tld := parser.Parse(c.domain)
			if sub != c.expectedSubdomain || root != c.expectedRootDomain || tld != c.expectedTopLevelDomain {
				t.Errorf("Parse(%s): expected %s, %s, %s, got %s, %s, %s", c.domain, c.expectedSubdomain, c.expectedRootDomain, c.expectedTopLevelDomain, sub, root, tld)
			}
		})
	}
}

func TestDomainParsingWithCustomTLDs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		domain                 string
		expectedSubdomain      string
		expectedRootDomain     string
		expectedTopLevelDomain string
	}{
		{"example.custom", "", "example", "custom"},
		{"subdomain.example.custom", "subdomain", "example", "custom"},
	}

	parser := hqgourl.NewDomainParser(
		hqgourl.DomainParserWithTLDs("custom"),
	)

	for _, c := range cases {
		c := c

		t.Run(fmt.Sprintf("Parse(%s)", c.domain), func(t *testing.T) {
			t.Parallel()

			sub, root, tld := parser.Parse(c.domain)
			if sub != c.expectedSubdomain || root != c.expectedRootDomain || tld != c.expectedTopLevelDomain {
				t.Errorf("Parse(%s): expected %s, %s, %s, got %s, %s, %s", c.domain, c.expectedSubdomain, c.expectedRootDomain, c.expectedTopLevelDomain, sub, root, tld)
			}
		})
	}
}

// func TestDomainParserEdgeCases(t *testing.T) {
// 	testCases := []struct {
// 		domain       string
// 		expectedSubdomain  string
// 		expectedRootDomain string
// 		expectedTopLevelDomain  string
// 	}{
// 		{"", "", "", ""},
// 		{"..", "", "", ""},
// 		{"noTLD", "", "noTLD", ""},
// 		{"123.456", "123", "456", ""},
// 		{"onlyroot", "", "onlyroot", ""},
// 		{"invalid..domain", "", "", ""},
// 		{"-invalid-.com", "", "", ""},
// 		{"sub.-domain.com", "sub", "", ""},
// 		{"sub.domain-.com", "sub", "domain-", "com"},
// 		{"localhost:8080", "", "", ""},
// 	}

// 	parser := hqgourl.NewDomainParser()

// 	for _, tc := range testCases {
// 		sub, root, tld := parser.Parse(tc.domain)
// 		if sub != tc.expectedSubdomain || root != tc.expectedRootDomain || tld != tc.expectedTopLevelDomain {
// 			t.Errorf("Parse(%s): expected %s, %s, %s, got %s, %s, %s", tc.domain, tc.expectedSubdomain, tc.expectedRootDomain, tc.expectedTopLevelDomain, sub, root, tld)
// 		}
// 	}
// }
