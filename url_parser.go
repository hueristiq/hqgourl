package hqgourl

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// URL extends the standard net/url URL struct with additional domain-related fields.
// It includes details like subdomain, root domain, and Top-Level Domain (TLD), along with
// standard URL components. This struct provides a comprehensive representation of a URL.
type URL struct {
	*url.URL // Embedding the standard URL struct for base functionalities.

	Subdomain      string // Subdomain of the URL.
	RootDomain     string // Root domain of the URL.
	TopLevelDomain string // Top-Level Domain (TLD) of the URL.
	Port           int    // Port number used in the URL.
	Extension      string // File extension derived from the URL path.
}

// URLParser encapsulates the logic for parsing URLs with additional domain-specific information.
// It enhances the standard URL parsing with the extraction of subdomain, root domain, and TLD.
// It also handles the addition of a default scheme if one is not present in the input URL.
type URLParser struct {
	scheme string // DefaultScheme is the default URL scheme to use if not specified in the URL.

	dp *DomainParser // DomainParser used for parsing the domain-specific details.
}

// WithDefaultScheme allows setting a default scheme for the URLParser.
// This default scheme is used if the input URL doesn't specify a scheme.
func (up *URLParser) WithDefaultScheme(scheme string) {
	URLParserWithDefaultScheme(scheme)(up)
}

// DefaultScheme returns the currently set default scheme of the URLParser.
func (up *URLParser) DefaultScheme() (scheme string) {
	return up.scheme
}

// Parse takes a raw URL string and parses it into a URL struct.
// It adds domain-specific details like subdomain, root domain, and TLD to the parsed URL.
// The method also ensures a default scheme is set if the URL does not specify one.
func (up *URLParser) Parse(rawURL string) (parsedURL *URL, err error) {
	parsedURL = &URL{}

	// Add default scheme if necessary
	if up.scheme != "" {
		rawURL = addScheme(rawURL, up.scheme)
	}

	// Standard URL parsing
	parsedURL.URL, err = url.Parse(rawURL)
	if err != nil {
		err = fmt.Errorf("error parsing URL: %w", err)

		return
	}

	// Split host and port, and handle errors
	parsedURL.Host, parsedURL.Port, err = splitHostPort(parsedURL.Host)
	if err != nil {
		err = fmt.Errorf("error splitting host and port: %w", err)

		return
	}

	domainRegex := regexp.MustCompile(`(?i)(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,}`)

	if domainRegex.MatchString(parsedURL.Host) {
		// Parse domain-specific parts
		parsedDomain := up.dp.Parse(parsedURL.Host)

		parsedURL.Subdomain, parsedURL.RootDomain, parsedURL.TopLevelDomain = parsedDomain.Sub, parsedDomain.Root, parsedDomain.TopLevel
	}

	// Extract file extension from the path
	parsedURL.Extension = path.Ext(parsedURL.Path)

	return
}

// URLParserOptionsFunc defines a function type for configuring a URLParser.
type URLParserOptionsFunc func(*URLParser)

// URLParserInterface defines the interface for URL parsing functionality.
type URLParserInterface interface {
	WithDefaultScheme(scheme string)

	DefaultScheme() (scheme string)

	Parse(rawURL string) (parsedURL *URL, err error)
}

var _ URLParserInterface = &URLParser{}

// NewURLParser creates a new URLParser with the given options.
// It initializes a DomainParser for parsing domain details and applies any additional configuration options.
func NewURLParser(opts ...URLParserOptionsFunc) (up *URLParser) {
	up = &URLParser{}

	// Initialize the DomainParser
	dp := NewDomainParser()
	up.dp = dp

	// Apply additional options
	for _, opt := range opts {
		opt(up)
	}

	return
}

// URLParserWithDefaultScheme returns a URLParserOptionsFunc to set a default scheme.
// This is useful when parsing URLs that may not have a scheme included.
func URLParserWithDefaultScheme(scheme string) URLParserOptionsFunc {
	return func(up *URLParser) {
		up.scheme = scheme
	}
}

// addScheme is a helper function that adds a scheme to the URL if it's missing.
// This ensures that the URL is parsed correctly as a network address rather than a relative path.
// This makes net/url.Parse() not put both host and path into the (relative) path.
func addScheme(inURL, scheme string) (outURL string) {
	switch {
	case strings.HasPrefix(inURL, "//"):
		outURL = scheme + ":" + inURL
	case strings.HasPrefix(inURL, "://"):
		outURL = scheme + inURL
	case !strings.Contains(inURL, "//"):
		outURL = scheme + "://" + inURL
	default:
		outURL = inURL
	}

	return
}

// splitHostPort separates the host and port in a network address.
// It is designed to handle both IPv4 and IPv6 addresses and gracefully manages URLs without a port.
// Unlike net.SplitHostPort(), it doesn't remove brackets from [IPv6] hosts.
func splitHostPort(address string) (host string, port int, err error) {
	host = address

	// Check for the last colon, which should separate host and port
	i := strings.LastIndex(address, ":")
	if i == -1 {
		return
	}

	// Handle IPv6 addresses enclosed in brackets
	if strings.HasPrefix(address, "[") && strings.Contains(address[i:], "]") {
		return
	}

	// Split the host and port
	host = address[:i]

	if port, err = strconv.Atoi(address[i+1:]); err != nil {
		return
	}

	return
}
