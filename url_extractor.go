package hqgourl

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/hueristiq/hqgourl/schemes"
	"github.com/hueristiq/hqgourl/tlds"
	"github.com/hueristiq/hqgourl/unicodes"
)

// URLExtractor is a struct that configures the URL extraction process.
// It allows specifying whether to include URL schemes and hosts in the extraction and supports
// custom regex patterns for these components.
type URLExtractor struct {
	withScheme        bool   // Indicates if the scheme part is mandatory in the URLs to be extracted.
	withSchemePattern string // Custom regex pattern for matching URL schemes, if provided.
	withHost          bool   // Indicates if the host part is mandatory in the URLs to be extracted.
	withHostPattern   string // Custom regex pattern for matching URL hosts, if provided.
}

// CompileRegex compiles a regex pattern based on the URLExtractor configuration.
// It dynamically constructs a regex pattern to accurately capture URLs from text,
// supporting various URL formats and components. The method ensures the regex captures
// the longest possible match for a URL, enhancing the accuracy of the extraction process.
func (e *URLExtractor) CompileRegex() (regex *regexp.Regexp) {
	schemePattern := URLExtractorSchemePattern

	if e.withScheme && e.withSchemePattern != "" {
		schemePattern = e.withSchemePattern
	}

	var asciiTLDs, unicodeTLDs []string

	for i, tld := range tlds.TLDs {
		if tld[0] >= utf8.RuneSelf {
			asciiTLDs = tlds.TLDs[:i:i]
			unicodeTLDs = tlds.TLDs[i:]

			break
		}
	}

	punycode := `xn--[a-z0-9-]+`
	knownTLDPattern := `(?:(?i)` + punycode + `|` + anyOf(append(asciiTLDs, tlds.PseudoTLDs...)...) + `\b|` + anyOf(unicodeTLDs...) + `)`
	domainPattern := `(?:` + _subdomainPattern + knownTLDPattern + `|localhost)`

	hostWithoutPortPattern := `(?:` + domainPattern + `|\[` + URLExtractorIPv6Pattern + `\]|\b` + URLExtractorIPv4Pattern + `\b)`
	hostWithPortOptionalPattern := `(?:` + hostWithoutPortPattern + URLExtractorPortOptionalPattern + `)`

	if e.withHost && e.withHostPattern != "" {
		hostWithPortOptionalPattern = e.withHostPattern
	}

	_IAuthorityPattern := `(?:` + _IUserInfoOptionalPattern + hostWithPortOptionalPattern + `)`
	_IAuthorityOptionalPattern := _IAuthorityPattern + `?`

	webURL := _IAuthorityPattern + `(?:/` + pathCont + `|/)?`

	// Emails pattern.
	email := `(?P<relaxedEmail>[a-zA-Z0-9._%\-+]+@` + hostWithPortOptionalPattern + `)`

	URLsWithSchemePattern := schemePattern + _IAuthorityOptionalPattern + pathCont

	if e.withHostPattern != "" {
		URLsWithSchemePattern = schemePattern + _IAuthorityPattern + `(?:/` + pathCont + `|/)?`
	}

	URLsWithHostPattern := webURL + `|` + email

	RelativeURLsPattern := `(\/[\w\/?=&#.-]*)|([\w\/?=&#.-]+?(?:\/[\w\/?=&#.-]+)+)`

	var pattern string

	switch {
	case e.withScheme:
		pattern = URLsWithSchemePattern
	case e.withHost:
		pattern = URLsWithSchemePattern + `|` + URLsWithHostPattern
	default:
		pattern = URLsWithSchemePattern + `|` + URLsWithHostPattern + `|` + RelativeURLsPattern
	}

	// Compiling the final regex pattern.
	regex = regexp.MustCompile(pattern)
	// Ensures the longest possible match is found.
	regex.Longest()

	return
}

// URLExtractorOptionsFunc defines a function type for configuring URLExtractor instances.
// This approach allows for flexible and fluent configuration of the extractor.
type URLExtractorOptionsFunc func(*URLExtractor)

// URLExtractorInterface defines the interface for URLExtractor, ensuring it implements certain methods.
type URLExtractorInterface interface {
	CompileRegex() (regex *regexp.Regexp)
}

const (
	_alphaCharacterSet          = `a-zA-Z`
	_digitCHaracterSet          = `0-9`
	_IUnreservedCharacterSet    = _alphaCharacterSet + _digitCHaracterSet + `\-\._~` + unicodes.AllowedUcsChar
	_IEndUnreservedCharacterSet = _alphaCharacterSet + _digitCHaracterSet + `\-_~` + unicodes.AllowedUcsCharMinusPunc
	_subDelimsCharacterSet      = `!\$&'\(\)\*\+,;=`
	_endSubDelimsCharacterSet   = `\$&\+=`
	_pctEncodingPattern         = `%[0-9a-fA-F]{2}`

	_IUserInfoPattern         = `(?:(?:[` + _IUnreservedCharacterSet + _subDelimsCharacterSet + `:]|` + _pctEncodingPattern + `)+@)`
	_IUserInfoOptionalPattern = _IUserInfoPattern + `?`

	URLExtractorIPv4Pattern         = `(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])`
	URLExtractorNonEmptyIPv6Pattern = `(?:` +
		// 7 colon-terminated chomps, followed by a final chomp or the rest of an elision.
		`(?:[0-9a-fA-F]{1,4}:){7}(?:[0-9a-fA-F]{1,4}|:)|` +
		// 6 chomps, followed by an IPv4 address or elision with final chomp or final elision.
		`(?:[0-9a-fA-F]{1,4}:){6}(?:` + URLExtractorIPv4Pattern + `|:[0-9a-fA-F]{1,4}|:)|` +
		// 5 chomps, followed by an elision with optional IPv4 or up to 2 final chomps.
		`(?:[0-9a-fA-F]{1,4}:){5}(?::` + URLExtractorIPv4Pattern + `|(?::[0-9a-fA-F]{1,4}){1,2}|:)|` +
		// 4 chomps, followed by an elision with optional IPv4 (optionally preceded by a chomp) or
		// up to 3 final chomps.
		`(?:[0-9a-fA-F]{1,4}:){4}(?:(?::[0-9a-fA-F]{1,4}){0,1}:` + URLExtractorIPv4Pattern + `|(?::[0-9a-fA-F]{1,4}){1,3}|:)|` +
		// 3 chomps, followed by an elision with optional IPv4 (preceded by up to 2 chomps) or
		// up to 4 final chomps.
		`(?:[0-9a-fA-F]{1,4}:){3}(?:(?::[0-9a-fA-F]{1,4}){0,2}:` + URLExtractorIPv4Pattern + `|(?::[0-9a-fA-F]{1,4}){1,4}|:)|` +
		// 2 chomps, followed by an elision with optional IPv4 (preceded by up to 3 chomps) or
		// up to 5 final chomps.
		`(?:[0-9a-fA-F]{1,4}:){2}(?:(?::[0-9a-fA-F]{1,4}){0,3}:` + URLExtractorIPv4Pattern + `|(?::[0-9a-fA-F]{1,4}){1,5}|:)|` +
		// 1 chomp, followed by an elision with optional IPv4 (preceded by up to 4 chomps) or
		// up to 6 final chomps.
		`(?:[0-9a-fA-F]{1,4}:){1}(?:(?::[0-9a-fA-F]{1,4}){0,4}:` + URLExtractorIPv4Pattern + `|(?::[0-9a-fA-F]{1,4}){1,6}|:)|` +
		// elision, followed by optional IPv4 (preceded by up to 5 chomps) or up to 7 final chomps.
		// `:` is an intentionally omitted alternative, to avoid matching `::`.
		`:(?:(?::[0-9a-fA-F]{1,4}){0,5}:` + URLExtractorIPv4Pattern + `|(?::[0-9a-fA-F]{1,4}){1,7})` +
		`)`
	URLExtractorIPv6Pattern = `(?:` + URLExtractorNonEmptyIPv6Pattern + `|::)`

	URLExtractorPortPattern         = `(?::[0-9]{1,4}|[1-5][0-9]{4}|6[0-5][0-9]{3}\b)`
	URLExtractorPortOptionalPattern = URLExtractorPortPattern + `?`

	midIPathSegmentChar = _IUnreservedCharacterSet + `%` + _subDelimsCharacterSet + `:@`
	endIPathSegmentChar = _IEndUnreservedCharacterSet + `%` + _endSubDelimsCharacterSet

	_IPrivateCharacters = `\x{E000}-\x{F8FF}\x{F0000}-\x{FFFFD}\x{100000}-\x{10FFFD}`

	midIChar  = `/?#\\` + midIPathSegmentChar + _IPrivateCharacters
	endIChar  = `/#` + endIPathSegmentChar + _IPrivateCharacters
	wellParen = `\((?:[` + midIChar + `]|\([` + midIChar + `]*\))*\)`
	wellBrack = `\[(?:[` + midIChar + `]|\[[` + midIChar + `]*\])*\]`
	wellBrace = `\{(?:[` + midIChar + `]|\{[` + midIChar + `]*\})*\}`
	wellAll   = wellParen + `|` + wellBrack + `|` + wellBrace
	pathCont  = `(?:[` + midIChar + `]*(?:` + wellAll + `|[` + endIChar + `]))+`

	_letter              = `\p{L}`
	_mark                = `\p{M}`
	_number              = `\p{N}`
	_IRICharctersPattern = `[` + _letter + _mark + _number + `](?:[` + _letter + _mark + _number + `\-]*[` + _letter + _mark + _number + `])?`

	_subdomainPattern = `(?:` + _IRICharctersPattern + `\.)+`
)

var (
	// URLExtractorSchemePattern defines a general pattern for matching URL schemes.
	// It matches any scheme that starts with alphabetical characters followed by any combination
	// of alphabets, dots, hyphens, or pluses, and ends with "://". It also matches any scheme
	// from a predefined list that does not require authority (host), ending with ":".
	URLExtractorSchemePattern = `(?:[a-zA-Z][a-zA-Z.\-+]*://|` + anyOf(schemes.SchemesNoAuthority...) + `:)`
	// URLExtractorKnownOfficialSchemePattern defines a pattern for matching officially recognized
	// URL schemes. This includes schemes like "http", "https", "ftp", etc., and is strictly based
	// on the schemes defined in the schemes.Schemes slice, ensuring a match ends with "://".
	URLExtractorKnownOfficialSchemePattern = `(?:` + anyOf(schemes.Schemes...) + `://)`
	// URLExtractorKnownUnofficialSchemePattern defines a pattern for matching unofficial or
	// less commonly used URL schemes. Similar to the official pattern but based on the
	// schemes.SchemesUnofficial slice, it supports schemes that might not be universally recognized
	// but are valid in specific contexts, ending with "://".
	URLExtractorKnownUnofficialSchemePattern = `(?:` + anyOf(schemes.SchemesUnofficial...) + `://)`
	// URLExtractorKnownNoAuthoritySchemePattern defines a pattern for matching schemes that
	// do not require an authority (host) component. This is useful for schemes like "mailto:",
	// "tel:", and others where a host is not applicable, ending with ":".
	URLExtractorKnownNoAuthoritySchemePattern = `(?:` + anyOf(schemes.SchemesNoAuthority...) + `:)`
	// URLExtractorKnownSchemePattern combines the patterns for officially recognized,
	// unofficial, and no-authority-required schemes into one comprehensive pattern. It is
	// case-insensitive (noted by "(?i)") and designed to match a wide range of schemes, accommodating
	// the broadest possible set of URLs.
	URLExtractorKnownSchemePattern = `(?:(?i)(?:` + anyOf(schemes.Schemes...) + `|` + anyOf(schemes.SchemesUnofficial...) + `)://|` + anyOf(schemes.SchemesNoAuthority...) + `:)`

	_ URLExtractorInterface = &URLExtractor{}
)

// NewURLExtractor creates a new URLExtractor instance with optional configuration.
// It applies the provided options to the extractor, allowing for customized behavior.
func NewURLExtractor(opts ...URLExtractorOptionsFunc) (extractor *URLExtractor) {
	extractor = &URLExtractor{}

	for _, opt := range opts {
		opt(extractor)
	}

	return
}

// URLExtractorWithScheme returns an option function to include URL schemes in the extraction process.
func URLExtractorWithScheme() URLExtractorOptionsFunc {
	return func(e *URLExtractor) {
		e.withScheme = true
	}
}

// URLExtractorWithSchemePattern returns an option function to specify a custom regex pattern
// for matching URL schemes. This allows for fine-tuned control over which schemes are considered valid.
func URLExtractorWithSchemePattern(pattern string) URLExtractorOptionsFunc {
	return func(e *URLExtractor) {
		e.withScheme = true
		e.withSchemePattern = pattern
	}
}

// URLExtractorWithHost returns an option function to include hosts in the URLs to be extracted.
// This can be used to ensure that only URLs with specified host components are captured.
func URLExtractorWithHost() URLExtractorOptionsFunc {
	return func(e *URLExtractor) {
		e.withHost = true
	}
}

// URLExtractorWithHostPattern returns an option function to specify a custom regex pattern
// for matching URL hosts. This is useful for targeting specific domain names or IP address formats.
func URLExtractorWithHostPattern(pattern string) URLExtractorOptionsFunc {
	return func(e *URLExtractor) {
		e.withHost = true
		e.withHostPattern = pattern
	}
}

// anyOf is a helper function that constructs a regex pattern for a set of strings.
// It simplifies the creation of regex patterns by automatically escaping and joining the provided strings.
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
