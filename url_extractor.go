package hqgourl

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/hueristiq/hqgourl/schemes"
	"github.com/hueristiq/hqgourl/tlds"
	"github.com/hueristiq/hqgourl/unicodes"
)

// URLExtractor is a struct for extracting URLs from text. It can be configured
// to recognize URLs based on specific schemes, hosts, and TLDs.
type URLExtractor struct {
	withScheme bool     // Determines if URL extraction is limited to certain schemes.
	schemes    []string // List of schemes to be considered in URL extraction.
	withHost   bool     // Flag to limit URL extraction to specific hosts.
	hosts      []string // List of hosts to be considered in URL extraction.
	withTlds   bool     // Flag to limit URL extraction to specific TLDs.
	tlds       []string // List of TLDs to be considered in URL extraction.
}

// WithScheme configures the URLExtractor to consider only URLs with the specified schemes.
func (e *URLExtractor) WithScheme(targetSchemes ...string) {
	URLExtractorWithScheme(targetSchemes...)(e)
}

// WithHost configures the URLExtractor to consider only URLs with the specified hosts.
func (e *URLExtractor) WithHost(targetHosts ...string) {
	URLExtractorWithHost(targetHosts...)(e)
}

// CompileRegex compiles a regular expression based on the configuration of the URLExtractor.
// This regex can be used to identify and extract URLs from text.
// The method combines various patterns to construct a comprehensive regex for URL extraction.
func (e *URLExtractor) CompileRegex() (regex *regexp.Regexp) {
	// 1. URLs with scheme
	URLsWithSchemePattern := `(?:(?i)(?:` + anyOf(schemes.Schemes...) + `|` + anyOf(schemes.SchemesUnofficial...) + `)://|` + anyOf(schemes.SchemesNoAuthority...) + `:)` + pathCont

	if len(e.schemes) > 0 {
		URLsWithSchemePattern = `(?i)(?:` + anyOf(e.schemes...) + `)://` + pathCont
	}

	// 2. URLs without scheme but has host
	var asciiTLDs, unicodeTLDs []string

	for i, tld := range tlds.TLDs {
		if tld[0] >= utf8.RuneSelf {
			asciiTLDs = tlds.TLDs[:i:i]
			unicodeTLDs = tlds.TLDs[i:]

			break
		}
	}

	punycode := `xn--[a-z0-9-]+`

	// Use \b to make sure ASCII TLDs are immediately followed by a word break.
	// We can't do that with unicode TLDs, as they don't see following
	// whitespace as a word break.
	TLDsPattern := `(?:(?i)` + punycode + `|` + anyOf(append(asciiTLDs, tlds.PseudoTLDs...)...) + `\b|` + anyOf(unicodeTLDs...) + `)`

	domain := _subdomainPattern + TLDsPattern

	// Web URLs pattern.
	hostName := `(?:` + domain + `|\[` + _IPv6AddressPattern + `\]|\b` + _IPv4AdressPattern + `\b)`

	webURL := hostName + _port + `(?:/` + pathCont + `|/)?`

	// Emails pattern.
	email := `(?P<relaxedEmail>[a-zA-Z0-9._%\-+]+@` + domain + `)`

	if len(e.hosts) > 0 {
		email = `(?P<relaxedEmail>[a-zA-Z0-9._%\-+]+@` + anyOf(e.hosts...) + `)`
	}

	URLsWithHostPattern := webURL + `|` + email + `|` + _nonEmptyIPv6AddressPattern

	// 3. URLs without scheme and host
	RelativeURLsPattern := `(\/[\w\/?=&#.-]*)|([\w\/?=&#.-]+?(?:\/[\w\/?=&#.-]+)+)`

	pattern := URLsWithSchemePattern + `|` + URLsWithHostPattern + `|` + RelativeURLsPattern

	if e.withScheme {
		pattern = URLsWithSchemePattern
	} else if e.withHost {
		pattern = URLsWithSchemePattern + `|` + URLsWithHostPattern
	}

	regex = regexp.MustCompile(pattern)

	regex.Longest()

	return
}

type URLExtractorOptionsFunc func(*URLExtractor)

type URLExtractorInterface interface {
	WithScheme(targetSchemes ...string)
	WithHost(targetHosts ...string)

	CompileRegex() (regex *regexp.Regexp)
}

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
	extractor = &URLExtractor{}

	// Apply configuration options to the extractor.
	for _, opt := range opts {
		opt(extractor)
	}

	return
}

// URLExtractorWithScheme returns a URLExtractorOptionsFunc that sets the target schemes for URL extraction.
// This function allows for limiting the URL extraction to specific URL schemes.
func URLExtractorWithScheme(targetSchemes ...string) URLExtractorOptionsFunc {
	return func(e *URLExtractor) {
		e.withScheme = true
		e.schemes = targetSchemes
	}
}

// URLExtractorWithHost returns a URLExtractorOptionsFunc that sets the target hosts for URL extraction.
// This function allows for limiting the URL extraction to URLs with specific hosts.
func URLExtractorWithHost(targetHosts ...string) URLExtractorOptionsFunc {
	return func(e *URLExtractor) {
		e.withHost = true
		e.schemes = targetHosts
	}
}

// URLExtractorWithTLD returns a URLExtractorOptionsFunc that sets the target TLDs for URL extraction.
// This function allows for limiting the URL extraction to URLs with specific TLDs.
func URLExtractorWithTLD(targetTLDs ...string) URLExtractorOptionsFunc {
	return func(e *URLExtractor) {
		e.withTlds = true
		e.tlds = targetTLDs
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
