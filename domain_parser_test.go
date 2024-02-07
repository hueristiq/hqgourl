package hqgourl_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hueristiq/hqgourl"
	"github.com/hueristiq/hqgourl/tlds"
)

func TestDomain_String(t *testing.T) {
	t.Parallel()

	cases := []struct {
		domainStruct         *hqgourl.Domain
		expectedDomainString string
	}{
		{
			&hqgourl.Domain{
				"",
				"example",
				"",
			},
			"example",
		},
		{
			&hqgourl.Domain{
				"",
				"example",
				"com",
			},
			"example.com",
		},
		{
			&hqgourl.Domain{
				"www",
				"example",
				"com",
			},
			"www.example.com",
		},
		{
			&hqgourl.Domain{
				"blog.www",
				"example",
				"com",
			},
			"blog.www.example.com",
		},
	}

	for _, c := range cases {
		c := c

		t.Run(fmt.Sprintf("Domain.String(%q)", c.domainStruct), func(t *testing.T) {
			t.Parallel()

			domainString := c.domainStruct.String()

			if domainString != c.expectedDomainString {
				t.Errorf("Domain.String(%q) = %v, want %v", c.domainStruct, domainString, c.expectedDomainString)
			}
		})
	}
}

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
		rawDomain            string
		expectedParsedDomain *hqgourl.Domain
	}{
		{
			"localhost",
			&hqgourl.Domain{
				Sub:      "",
				Root:     "localhost",
				TopLevel: "",
			},
		},
		{
			"example.com",
			&hqgourl.Domain{
				Sub:      "",
				Root:     "example",
				TopLevel: "com",
			},
		},
		{
			"www.example.com",
			&hqgourl.Domain{
				Sub:      "www",
				Root:     "example",
				TopLevel: "com",
			},
		},
		{
			"example.co.uk",
			&hqgourl.Domain{
				Sub:      "",
				Root:     "example",
				TopLevel: "co.uk",
			},
		},
		{
			"www.example.co.uk",
			&hqgourl.Domain{
				Sub:      "www",
				Root:     "example",
				TopLevel: "co.uk",
			},
		},
	}

	dp := hqgourl.NewDomainParser()

	for _, c := range cases {
		c := c

		t.Run(fmt.Sprintf("Parse(%q)", c.rawDomain), func(t *testing.T) {
			t.Parallel()

			parsedDomain := dp.Parse(c.rawDomain)

			if !reflect.DeepEqual(parsedDomain, c.expectedParsedDomain) {
				t.Errorf("Parse(%q) = %+v, want %+v", c.rawDomain, parsedDomain, c.expectedParsedDomain)
			}
		})
	}
}

func TestDomainParsingWithCustomTLDs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		rawDomain            string
		expectedParsedDomain *hqgourl.Domain
	}{
		{
			"example.custom",
			&hqgourl.Domain{
				Sub:      "",
				Root:     "example",
				TopLevel: "custom",
			},
		},
		{
			"subdomain.example.custom",
			&hqgourl.Domain{
				Sub:      "subdomain",
				Root:     "example",
				TopLevel: "custom",
			},
		},
	}

	dp := hqgourl.NewDomainParser(
		hqgourl.DomainParserWithTLDs("custom"),
	)

	for _, c := range cases {
		c := c

		t.Run(fmt.Sprintf("Parse(%q)", c.rawDomain), func(t *testing.T) {
			t.Parallel()

			parsedDomain := dp.Parse(c.rawDomain)

			if !reflect.DeepEqual(parsedDomain, c.expectedParsedDomain) {
				t.Errorf("Parse(%q) = %+v, want %+v", c.rawDomain, parsedDomain, c.expectedParsedDomain)
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

		parsedDomain := dp.Parse(domain)
		if parsedDomain.TopLevel != TLD {
			t.Errorf("Parse(%q) = %q, %q, %q; want %q, %q, %q", domain, "", "example", TLD, "", "example", parsedDomain.TopLevel)
		}
	}
}
