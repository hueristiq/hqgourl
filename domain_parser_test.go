package hqgourl_test

import (
	"testing"

	"github.com/hueristiq/hqgourl"
)

func TestDomainParsing(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		domain       string
		expectedSub  string
		expectedRoot string
		expectedTLD  string
	}{
		{"www.example.com", "www", "example", "com"},
		{"subdomain.example.co.uk", "subdomain", "example", "co.uk"},
		{"example.com", "", "example", "com"},
		{"localhost", "", "localhost", ""},
	}

	parser := hqgourl.NewDomainParser()

	for _, tc := range testCases {
		sub, root, tld := parser.Parse(tc.domain)
		if sub != tc.expectedSub || root != tc.expectedRoot || tld != tc.expectedTLD {
			t.Errorf("Parse(%s): expected %s, %s, %s, got %s, %s, %s", tc.domain, tc.expectedSub, tc.expectedRoot, tc.expectedTLD, sub, root, tld)
		}
	}
}

func TestDomainParsingWithCustomTLDs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		domain       string
		expectedSub  string
		expectedRoot string
		expectedTLD  string
	}{
		{"example.custom", "", "example", "custom"},
		{"subdomain.example.custom", "subdomain", "example", "custom"},
	}

	parser := hqgourl.NewDomainParser(
		hqgourl.DomainParserWithTLDs("custom"),
	)

	for _, tc := range testCases {
		sub, root, tld := parser.Parse(tc.domain)
		if sub != tc.expectedSub || root != tc.expectedRoot || tld != tc.expectedTLD {
			t.Errorf("Parse(%s): expected %s, %s, %s, got %s, %s, %s", tc.domain, tc.expectedSub, tc.expectedRoot, tc.expectedTLD, sub, root, tld)
		}
	}
}

// func TestDomainParserEdgeCases(t *testing.T) {
// 	testCases := []struct {
// 		domain       string
// 		expectedSub  string
// 		expectedRoot string
// 		expectedTLD  string
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
// 		if sub != tc.expectedSub || root != tc.expectedRoot || tld != tc.expectedTLD {
// 			t.Errorf("Parse(%s): expected %s, %s, %s, got %s, %s, %s", tc.domain, tc.expectedSub, tc.expectedRoot, tc.expectedTLD, sub, root, tld)
// 		}
// 	}
// }
