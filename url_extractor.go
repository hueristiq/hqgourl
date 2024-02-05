package hqgourl

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/hueristiq/hqgourl/schemes"
	"github.com/hueristiq/hqgourl/tlds"
	"github.com/hueristiq/hqgourl/unicodes"
)

// URLExtractorStrictness defines the strictness levels for URL extraction.
type URLExtractorStrictness int

// IsHigh checks if the strictness level is high.
func (s URLExtractorStrictness) IsHigh() bool {
	return s == URLExtractorHighStrictness
}

// IsMedium checks if the strictness level is medium.
func (s URLExtractorStrictness) IsMedium() bool {
	return s == URLExtractorMediumStrictness
}

// IsLow checks if the strictness level is low.
func (s URLExtractorStrictness) IsLow() bool {
	return s == URLExtractorLowStrictness
}

// URLExtractor is a struct that extracts URLs from text with configurable strictness.
type URLExtractor struct {
	strictness URLExtractorStrictness // The strictness level for URL extraction.
}

// Strictness returns the current strictness level of the URLExtractor.
func (e *URLExtractor) Strictness() (strictness URLExtractorStrictness) {
	return e.strictness
}

// CompileRegex compiles the regex based on the URLExtractor's strictness level.
// It combines schemes, TLDs, and path patterns to create a regex that matches URLs.
func (e *URLExtractor) CompileRegex() (regex *regexp.Regexp) {
	// Scheme Pattern: Matches both official and unofficial schemes.
	schemePattern := `(?:(?i)(?:` + anyOf(schemes.Schemes...) + `|` + anyOf(schemes.SchemesUnofficial...) + `)://|` + anyOf(schemes.SchemesNoAuthority...) + `:)`

	// TLD Pattern: Separates ASCII and Unicode TLDs, includes pseudo TLDs and punycode representation.
	var asciiTLDs, unicodeTLDs []string

	for i, tld := range tlds.TLDs {
		if tld[0] >= utf8.RuneSelf {
			asciiTLDs = tlds.TLDs[:i:i]
			unicodeTLDs = tlds.TLDs[i:]

			break
		}
	}

	punycode := `xn--[a-z0-9-]+`

	TLDsPattern := `(?:(?i)` + punycode + `|` + anyOf(append(asciiTLDs, tlds.PseudoTLDs...)...) + `\b|` + anyOf(unicodeTLDs...) + `)`

	// Domain and Host Patterns: Include support for subdomains, IPv4, and IPv6 addresses.
	domainPattern := _subdomainPattern + TLDsPattern
	hostPattern := `(?:` + domainPattern + `|\[` + _IPv6AddressPattern + `\]|\b` + _IPv4AdressPattern + `\b)`

	webURL := hostPattern + _port + `(?:/` + pathCont + `|/)?`

	// Emails pattern.
	email := `(?P<relaxedEmail>[a-zA-Z0-9._%\-+]+@` + domainPattern + `)`

	URLsWithSchemePattern := schemePattern + pathCont
	URLsWithHostPattern := webURL + `|` + email + `|` + _nonEmptyIPv6AddressPattern
	RelativeURLsPattern := `(\/[\w\/?=&#.-]*)|([\w\/?=&#.-]+?(?:\/[\w\/?=&#.-]+)+)`

	// Selecting the pattern based on the specified strictness level.
	pattern := ``

	switch {
	case e.strictness.IsLow():
		pattern = URLsWithSchemePattern + `|` + URLsWithHostPattern + `|` + RelativeURLsPattern
	case e.strictness.IsMedium():
		pattern = URLsWithSchemePattern + `|` + URLsWithHostPattern
	case e.strictness.IsHigh():
		pattern = URLsWithSchemePattern
	}

	// Compiling the final regex pattern.
	regex = regexp.MustCompile(pattern)
	// Ensures the longest possible match is found.
	regex.Longest()

	return
}

// URLExtractorOptionsFunc defines a function type for configuring URLExtractor instances.
type URLExtractorOptionsFunc func(*URLExtractor)

// URLExtractorInterface defines the interface for URLExtractor, ensuring it implements certain methods.
type URLExtractorInterface interface {
	Strictness() (strictness URLExtractorStrictness)
	CompileRegex() (regex *regexp.Regexp)
}

// Predefined strictness levels.
const (
	URLExtractorLowStrictness    = URLExtractorStrictness(1)
	URLExtractorMediumStrictness = URLExtractorStrictness(2)
	URLExtractorHighStrictness   = URLExtractorStrictness(3)
)

const (
	// pathCont is based on https://www.rfc-editor.org/rfc/rfc3987#section-2.2
	// but does not match separators anywhere or most punctuation in final position,
	// to avoid creating asymmetries like
	// `Did you know that **<a href="...">https://example.com/**</a> is reserved for documentation?`
	// from `Did you know that **https://example.com/** is reserved for documentation?`.
	unreservedChar      = `a-zA-Z0-9\-\._~`
	endUnreservedChar   = `a-zA-Z0-9\-_~`
	midSubDelimChar     = `!\$&'\*\+,;=`
	endSubDelimChar     = `\$&\+=`
	midIPathSegmentChar = unreservedChar + `%` + midSubDelimChar + `:@` + unicodes.AllowedUcsChar
	endIPathSegmentChar = endUnreservedChar + `%` + endSubDelimChar + unicodes.AllowedUcsCharMinusPunc
	iPrivateChar        = `\x{E000}-\x{F8FF}\x{F0000}-\x{FFFFD}\x{100000}-\x{10FFFD}`
	midIChar            = `/?#\\` + midIPathSegmentChar + iPrivateChar
	endIChar            = `/#` + endIPathSegmentChar + iPrivateChar
	wellParen           = `\((?:[` + midIChar + `]|\([` + midIChar + `]*\))*\)`
	wellBrack           = `\[(?:[` + midIChar + `]|\[[` + midIChar + `]*\])*\]`
	wellBrace           = `\{(?:[` + midIChar + `]|\{[` + midIChar + `]*\})*\}`
	wellAll             = wellParen + `|` + wellBrack + `|` + wellBrace
	pathCont            = `(?:[` + midIChar + `]*(?:` + wellAll + `|[` + endIChar + `]))+`

	_letter              = `\p{L}`
	_mark                = `\p{M}`
	_number              = `\p{N}`
	_IRICharctersPattern = `[` + _letter + _mark + _number + `](?:[` + _letter + _mark + _number + `\-]*[` + _letter + _mark + _number + `])?`

	_subdomainPattern = `(?:` + _IRICharctersPattern + `\.)+`

	_IPv4AdressPattern          = `(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])`
	_nonEmptyIPv6AddressPattern = `(?:` +
		// 7 colon-terminated chomps, followed by a final chomp or the rest of an elision.
		`(?:[0-9a-fA-F]{1,4}:){7}(?:[0-9a-fA-F]{1,4}|:)|` +
		// 6 chomps, followed by an IPv4 address or elision with final chomp or final elision.
		`(?:[0-9a-fA-F]{1,4}:){6}(?:` + _IPv4AdressPattern + `|:[0-9a-fA-F]{1,4}|:)|` +
		// 5 chomps, followed by an elision with optional IPv4 or up to 2 final chomps.
		`(?:[0-9a-fA-F]{1,4}:){5}(?::` + _IPv4AdressPattern + `|(?::[0-9a-fA-F]{1,4}){1,2}|:)|` +
		// 4 chomps, followed by an elision with optional IPv4 (optionally preceded by a chomp) or
		// up to 3 final chomps.
		`(?:[0-9a-fA-F]{1,4}:){4}(?:(?::[0-9a-fA-F]{1,4}){0,1}:` + _IPv4AdressPattern + `|(?::[0-9a-fA-F]{1,4}){1,3}|:)|` +
		// 3 chomps, followed by an elision with optional IPv4 (preceded by up to 2 chomps) or
		// up to 4 final chomps.
		`(?:[0-9a-fA-F]{1,4}:){3}(?:(?::[0-9a-fA-F]{1,4}){0,2}:` + _IPv4AdressPattern + `|(?::[0-9a-fA-F]{1,4}){1,4}|:)|` +
		// 2 chomps, followed by an elision with optional IPv4 (preceded by up to 3 chomps) or
		// up to 5 final chomps.
		`(?:[0-9a-fA-F]{1,4}:){2}(?:(?::[0-9a-fA-F]{1,4}){0,3}:` + _IPv4AdressPattern + `|(?::[0-9a-fA-F]{1,4}){1,5}|:)|` +
		// 1 chomp, followed by an elision with optional IPv4 (preceded by up to 4 chomps) or
		// up to 6 final chomps.
		`(?:[0-9a-fA-F]{1,4}:){1}(?:(?::[0-9a-fA-F]{1,4}){0,4}:` + _IPv4AdressPattern + `|(?::[0-9a-fA-F]{1,4}){1,6}|:)|` +
		// elision, followed by optional IPv4 (preceded by up to 5 chomps) or up to 7 final chomps.
		// `:` is an intentionally omitted alternative, to avoid matching `::`.
		`:(?:(?::[0-9a-fA-F]{1,4}){0,5}:` + _IPv4AdressPattern + `|(?::[0-9a-fA-F]{1,4}){1,7})` +
		`)`
	_IPv6AddressPattern = `(?:` + _nonEmptyIPv6AddressPattern + `|::)`
	_port               = `(?::[0-9]+)?`
)

var _ URLExtractorInterface = &URLExtractor{}

// NewURLExtractor creates a new instance of URLExtractor with the provided options.
// This function allows for flexible configuration of the URLExtractor.
func NewURLExtractor(opts ...URLExtractorOptionsFunc) (extractor *URLExtractor) {
	extractor = &URLExtractor{
		strictness: URLExtractorLowStrictness,
	}

	for _, opt := range opts {
		opt(extractor)
	}

	return
}

// URLExtractorWithHighStrictness returns an option function to set high strictness.
func URLExtractorWithHighStrictness() URLExtractorOptionsFunc {
	return func(e *URLExtractor) {
		e.strictness = URLExtractorHighStrictness
	}
}

// URLExtractorWithMediumStrictness returns an option function to set medium strictness.
func URLExtractorWithMediumStrictness() URLExtractorOptionsFunc {
	return func(e *URLExtractor) {
		e.strictness = URLExtractorMediumStrictness
	}
}

// URLExtractorWithLowStrictness returns an option function to set low strictness.
func URLExtractorWithLowStrictness() URLExtractorOptionsFunc {
	return func(e *URLExtractor) {
		e.strictness = URLExtractorLowStrictness
	}
}

// anyOf is a helper function that constructs a regex pattern for a set of strings.
// It is used in compiling the regex patterns for URL extraction.
func anyOf(strs ...string) string {
	var b strings.Builder

	b.WriteString("(?:")

	for i, s := range strs {
		if i != 0 {
			b.WriteByte('|')
		}

		b.WriteString(regexp.QuoteMeta(s))
	}

	b.WriteByte(')')

	return b.String()
}
