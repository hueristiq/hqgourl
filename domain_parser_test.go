package hqgourl_test

import (
	"fmt"
	"testing"

	"github.com/hueristiq/hqgourl"
	"github.com/hueristiq/hqgourl/tlds"
)

func TestNewDomainParser(t *testing.T) {
	t.Parallel()

	dp := hqgourl.NewDomainParser()
	if dp == nil {
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

	dp := hqgourl.NewDomainParser()

	for _, c := range cases {
		c := c

		t.Run(fmt.Sprintf("Parse(%s)", c.domain), func(t *testing.T) {
			t.Parallel()

			sub, root, tld := dp.Parse(c.domain)
			if sub != c.expectedSubdomain || root != c.expectedRootDomain || tld != c.expectedTopLevelDomain {
				t.Errorf("Parse(%q) = %q, %q, %q; want %q, %q, %q", c.domain, sub, root, tld, c.expectedSubdomain, c.expectedRootDomain, c.expectedTopLevelDomain)
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

	dp := hqgourl.NewDomainParser(
		hqgourl.DomainParserWithTLDs("custom"),
	)

	for _, c := range cases {
		c := c

		t.Run(fmt.Sprintf("Parse(%s)", c.domain), func(t *testing.T) {
			t.Parallel()

			sub, root, tld := dp.Parse(c.domain)
			if sub != c.expectedSubdomain || root != c.expectedRootDomain || tld != c.expectedTopLevelDomain {
				t.Errorf("Parse(%q) = %q, %q, %q; want %q, %q, %q", c.domain, sub, root, tld, c.expectedSubdomain, c.expectedRootDomain, c.expectedTopLevelDomain)
			}
		})
	}
}

func TestDomainParserWithStandardAndPseudoTLDs(t *testing.T) {
	t.Parallel()

	dp := hqgourl.NewDomainParser()

	TLDs := []string{}

	TLDs = append(TLDs, tlds.TLDs...)
	TLDs = append(TLDs, tlds.PseudoTLDs...)

	for _, TLD := range TLDs {
		domain := "example." + TLD

		_, _, expectedTLD := dp.Parse(domain)
		if expectedTLD != TLD {
			t.Errorf("Parse(%q) = %q, %q, %q; want %q, %q, %q", domain, "", "example", TLD, "", "example", expectedTLD)
		}
	}
}
