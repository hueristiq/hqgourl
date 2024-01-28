package hqgourl

import (
	"index/suffixarray"
	"strings"

	"github.com/hueristiq/hqgourl/tlds"
)

// DomainParser encapsulates the logic for parsing domain names.
// It uses a suffix array to efficiently identify and extract components like subdomain,
// root domain, and TLD (Top-Level Domain) from a given domain string.
type DomainParser struct {
	sa *suffixarray.Index
}

// Parse takes a domain string and splits it into its subdomain, root domain, and TLD components.
func (dp *DomainParser) Parse(domain string) (subdomain, rootDomain, TLD string) {
	// Split the domain into parts based on '.'
	parts := strings.Split(domain, ".")

	// Find the offset index for the TLD in the domain parts
	TLDOffset := dp.findTLDOffset(parts)

	// Extract subdomain, root domain, and TLD based on the TLD offset
	subdomain = strings.Join(parts[:TLDOffset], ".")
	rootDomain = parts[TLDOffset]
	TLD = strings.Join(parts[TLDOffset+1:], ".")

	return
}

// findTLDOffset identifies the index at which the TLD starts in a split domain string.
// The method uses a suffix array to perform an efficient search for known TLDs within the domain parts.
// It iterates in reverse to correctly handle cases where subdomains might resemble TLDs.
func (dp *DomainParser) findTLDOffset(parts []string) (offset int) {
	for i := 2; i > 0; i-- {
		startPoint := len(parts) - i
		if startPoint < 0 {
			break
		}

		TLDPart := strings.Join(parts[startPoint:], ".")

		indicies := dp.sa.Lookup([]byte(TLDPart), -1)

		if len(indicies) > 0 {
			offset = (len(parts) - (i + 1))

			if offset >= 0 {
				return
			}
		}
	}

	return
}

// DomainParserOptionsFunc is a function type that applies configuration options to a DomainParser.
// This allows for flexible initialization of the parser with specific settings or custom data.
type DomainParserOptionsFunc func(*DomainParser)

// DomainParserInterface defines an interface for a domain parser, ensuring consistency
// and standardization in how domain parsing functionality is implemented and exposed.
type DomainParserInterface interface {
	Parse(domain string) (subdomain, rootDomain, TLD string)
	findTLDOffset(parts []string) (offset int)
}

var _ DomainParserInterface = &DomainParser{}

// NewDomainParser initializes a new instance of DomainParser.
// It preloads the suffix array with a comprehensive list of TLDs and pseudo-TLDs for accurate parsing.
// This function can be extended with additional options to customize the parser's behavior.
func NewDomainParser(opts ...DomainParserOptionsFunc) (dp *DomainParser) {
	dp = &DomainParser{}

	// Combine standard and pseudo-TLDs
	TLDs := []string{}

	TLDs = append(TLDs, tlds.TLDs...)
	TLDs = append(TLDs, tlds.PseudoTLDs...)

	data := []byte("\x00" + strings.Join(TLDs, "\x00") + "\x00")

	// Initialize the suffix array with the combined TLDs
	dp.sa = suffixarray.New(data)

	// Apply any additional options
	for _, opt := range opts {
		opt(dp)
	}

	return
}

// DomainParserWithTLDs is an option function that allows for the initialization of a DomainParser
// with a custom set of TLDs. This is useful in scenarios where non-standard or niche TLDs are in use,
// and accurate domain parsing is required for these specific cases.
func DomainParserWithTLDs(TLDs ...string) DomainParserOptionsFunc {
	data := []byte("\x00" + strings.Join(TLDs, "\x00") + "\x00")

	sa := suffixarray.New(data)

	return func(dp *DomainParser) {
		dp.sa = sa
	}
}
