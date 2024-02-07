package hqgourl

import (
	"index/suffixarray"
	"strings"

	"github.com/hueristiq/hqgourl/tlds"
)

// DomainParser encapsulates the logic for parsing domain names into their constituent parts:
// subdomains, root domains, and top-level domains (TLDs). It leverages a suffix array for efficient
// search and extraction of these components from a full domain string.
type DomainParser struct {
	sa *suffixarray.Index // Suffix array optimized for quick TLD lookup.
}

// Parse takes a full domain string and splits it into its constituent parts: subdomain,
// root domain, and TLD. This method efficiently identifies the TLD using a suffix array
// and separates the remaining parts of the domain accordingly.
func (dp *DomainParser) Parse(domain string) (subdomain, rootDomain, TLD string) {
	// Split the domain into parts based on '.'
	parts := strings.Split(domain, ".")

	if len(parts) <= 1 {
		rootDomain = domain

		return
	}

	// Identify the index where the TLD begins using the findTLDOffset method.
	TLDOffset := dp.findTLDOffset(parts)

	// Based on the TLD offset, separate the domain string into subdomain, root domain, and TLD.
	subdomain = strings.Join(parts[:TLDOffset], ".")
	rootDomain = parts[TLDOffset]
	TLD = strings.Join(parts[TLDOffset+1:], ".")

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

// DomainParserOptionsFunc is a function type designed for configuring a DomainParser instance.
// It allows for the application of customization options, such as specifying custom TLDs.
type DomainParserOptionsFunc func(*DomainParser)

// DomainParserInterface ensures that any domain parser implementation provides a standard
// method set for parsing domain names, promoting consistency and reliability in usage.
type DomainParserInterface interface {
	Parse(domain string) (subdomain, rootDomain, TLD string)

	findTLDOffset(parts []string) (offset int)
}

var _ DomainParserInterface = &DomainParser{} // Ensures DomainParser implements DomainParserInterface.

// NewDomainParser creates and initializes a DomainParser with a comprehensive list of TLDs,
// including both standard and pseudo-TLDs. This setup ensures accurate parsing across a wide
// range of domain names. Additional options can be applied to customize the parser further.
func NewDomainParser(opts ...DomainParserOptionsFunc) (dp *DomainParser) {
	dp = &DomainParser{}

	// Combine standard and pseudo-TLDs for comprehensive coverage.
	TLDs := []string{}

	TLDs = append(TLDs, tlds.TLDs...)
	TLDs = append(TLDs, tlds.PseudoTLDs...)

	data := []byte("\x00" + strings.Join(TLDs, "\x00") + "\x00")

	// Initialize the suffix array with TLD data.
	dp.sa = suffixarray.New(data)

	// Apply any additional options
	for _, opt := range opts {
		opt(dp)
	}

	return
}

// DomainParserWithTLDs allows for the initialization of the DomainParser with a custom set of TLDs.
// This is particularly useful for applications requiring parsing of non-standard or niche TLDs.
func DomainParserWithTLDs(TLDs ...string) DomainParserOptionsFunc {
	data := []byte("\x00" + strings.Join(TLDs, "\x00") + "\x00")

	sa := suffixarray.New(data)

	return func(dp *DomainParser) {
		dp.sa = sa // Override the suffix array with custom TLD data.
	}
}
