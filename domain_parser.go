package hqgourl

import (
	"index/suffixarray"
	"strings"

	"github.com/hueristiq/hqgourl/tlds"
)

// Domain struct represents the structure of a parsed domain name, including its subdomain, root domain, and top-level domain (TLD).
type Domain struct {
	Sub      string
	Root     string
	TopLevel string
}

// String assembles the domain components back into a full domain string.
func (d *Domain) String() (domain string) {
	var parts []string

	if d.Sub != "" {
		parts = append(parts, d.Sub)
	}

	parts = append(parts, d.Root)

	if d.TopLevel != "" {
		parts = append(parts, d.TopLevel)
	}

	domain = strings.Join(parts, ".")

	return
}

// DomainInterface defines a standard interface for any domain representation.
type DomainInterface interface {
	String() (domain string)
}

// DomainParser encapsulates the logic for parsing full domain strings into their constituent parts:
// subdomains, root domains, and top-level domains (TLDs). It leverages a suffix array for efficient
// search and extraction of these components from a full domain string.
type DomainParser struct {
	sa *suffixarray.Index
}

// Parse takes a full domain string and splits it into its constituent parts: subdomain,
// root domain, and TLD. This method efficiently identifies the TLD using a suffix array
// and separates the remaining parts of the domain accordingly.
func (dp *DomainParser) Parse(domain string) (parsedDomain *Domain) {
	parsedDomain = &Domain{}

	// Split the domain into parts based on '.'
	parts := strings.Split(domain, ".")

	if len(parts) <= 1 {
		parsedDomain.Root = domain

		return
	}

	// Identify the index where the TLD begins using the findTLDOffset method.
	TLDOffset := dp.findTLDOffset(parts)

	// Based on the TLD offset, separate the domain string into subdomain, root domain, and TLD.
	parsedDomain.Sub = strings.Join(parts[:TLDOffset], ".")
	parsedDomain.Root = parts[TLDOffset]
	parsedDomain.TopLevel = strings.Join(parts[TLDOffset+1:], ".")

	return
}

// findTLDOffset determines the starting index of the TLD within a domain split into parts.
// It reverses through the parts of the domain to accurately handle cases where subdomains may
// mimic TLDs. The method uses the suffix array to find known TLDs efficiently.
func (dp *DomainParser) findTLDOffset(parts []string) (offset int) {
	partsLength := len(parts)
	partsLastIndex := partsLength - 1

	for i := partsLastIndex; i >= 0; i-- {
		// Construct a potential TLD from the current part to the end.
		TLD := strings.Join(parts[i:], ".")

		// Search for the TLD in the suffix array.
		indices := dp.sa.Lookup([]byte(TLD), -1)

		// If a match is found, update the offset, else break.
		if len(indices) > 0 {
			offset = i - 1
		} else {
			break
		}
	}

	return
}

// DomainParserInterface defines a standard interface for any DomainParser representation.
type DomainParserInterface interface {
	Parse(domain string) (parsedDomain *Domain)
	findTLDOffset(parts []string) (offset int)
}

// DomainParserOptionsFunc is a function type designed for configuring a DomainParser instance.
// It allows for the application of customization options, such as specifying custom TLDs.
type DomainParserOptionsFunc func(*DomainParser)

// Ensure type compatibility with interfaces.
var (
	_ DomainInterface       = &Domain{}
	_ DomainParserInterface = &DomainParser{}
)

// NewDomainParser creates and initializes a DomainParser with a comprehensive list of TLDs,
// including both standard and pseudo-TLDs. This setup ensures accurate parsing across a wide
// range of domain names. Additional options can be applied to customize the parser further.
func NewDomainParser(opts ...DomainParserOptionsFunc) (dp *DomainParser) {
	dp = &DomainParser{}

	// Combine standard and pseudo-TLDs for comprehensive coverage.
	TLDs := []string{}

	TLDs = append(TLDs, tlds.TLDs...)
	TLDs = append(TLDs, tlds.PseudoTLDs...)

	// Initialize the suffix array with TLD data.
	dp.sa = suffixarray.New([]byte("\x00" + strings.Join(TLDs, "\x00") + "\x00"))

	// Apply any additional options
	for _, opt := range opts {
		opt(dp)
	}

	return
}

// DomainParserWithTLDs allows for the initialization of the DomainParser with a custom set of TLDs.
// This is particularly useful for applications requiring parsing of non-standard or niche TLDs.
func DomainParserWithTLDs(TLDs ...string) DomainParserOptionsFunc {
	return func(dp *DomainParser) {
		dp.sa = suffixarray.New([]byte("\x00" + strings.Join(TLDs, "\x00") + "\x00"))
	}
}
