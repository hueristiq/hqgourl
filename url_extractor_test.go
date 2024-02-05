package hqgourl_test

import (
	"testing"

	"github.com/hueristiq/hqgourl"
)

func TestNewExtractor(t *testing.T) {
	t.Parallel()

	extractor := hqgourl.NewURLExtractor()

	if extractor == nil {
		t.Error("NewURLExtractor() = nil; want non-nil")
	}

	if extractor.Strictness() != hqgourl.URLExtractorLowStrictness {
		t.Errorf("NewURLExtractor().Strictness() = %v; want %v", extractor.Strictness(), hqgourl.URLExtractorLowStrictness)
	}

	highStrictnessExtractor := hqgourl.NewURLExtractor(hqgourl.URLExtractorWithHighStrictness())

	if highStrictnessExtractor.Strictness() != hqgourl.URLExtractorHighStrictness {
		t.Errorf("NewURLExtractor(URLExtractorWithHighStrictness()).Strictness() = %v; want %v", highStrictnessExtractor.Strictness(), hqgourl.URLExtractorHighStrictness)
	}

	mediumStrictnessExtractor := hqgourl.NewURLExtractor(hqgourl.URLExtractorWithMediumStrictness())

	if mediumStrictnessExtractor.Strictness() != hqgourl.URLExtractorMediumStrictness {
		t.Errorf("NewURLExtractor(URLExtractorWithMediumStrictness()).Strictness() = %v; want %v", mediumStrictnessExtractor.Strictness(), hqgourl.URLExtractorMediumStrictness)
	}

	lowStrictnessExtractor := hqgourl.NewURLExtractor(hqgourl.URLExtractorWithLowStrictness())

	if lowStrictnessExtractor.Strictness() != hqgourl.URLExtractorLowStrictness {
		t.Errorf("NewURLExtractor(URLExtractorWithLowStrictness()).Strictness() = %v; want %v", lowStrictnessExtractor.Strictness(), hqgourl.URLExtractorLowStrictness)
	}
}

func TestCompileRegex(t *testing.T) {
	t.Parallel()

	extractor := hqgourl.NewURLExtractor()

	regex := extractor.CompileRegex()

	if regex == nil {
		t.Errorf("CompileRegex() = nil; want non-nil")
	}
}

func TestURLExtractionWithLowStrictness(t *testing.T) {
	t.Parallel()

	extractor := hqgourl.NewURLExtractor()

	regex := extractor.CompileRegex()

	testCases := []struct {
		text string
		want []string
	}{
		{
			`
			Localhost URL: http://localhost
			Localhost URL with Port: http://localhost:8000/home
			Standard URL: https://www.example.com
			URL with Port: http://www.example.com:8080/page
			URL with HTTPS and Query: https://www.example.com/search?q=openai
			URL with Path: https://www.example.com/resources/docs/guide
			URL with Fragment: https://www.example.com/about#team
			URL with User Info: https://user:password@example.com
			Complex URL: https://www.example.com:8080/search?q=openai#results
			International URL: https://www.例子.公司.cn
			URL with IPv4 Address: http://192.168.1.1/setup
			URL with IPv4 and Port: http://192.168.1.1:8080/setup
			URL with IPv6 Address: http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]
			URL with IPv6 and Port: http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080
			URL with Special Characters: https://example.com/products?id=50&category=books
			URL with Encoded Characters: https://example.com/search?name=John%20Doe&age=30
			URL with Multiple Subdomains: https://subdomain.example.co.uk
			File Protocol URL: file:///C:/path/to/file.txt
			Mailto Protocol URL: mailto:user@example.com
			FTP URL with Port: ftp://user:password@ftp.example.com:21
			Data URL: data:text/plain;base64,aGVsbG8=
			URL with JavaScript Protocol: javascript:alert('Hello World');
			URL with Spaces: https://www.example.com/test space
			URL in Text: Check out this link: https://www.example.com, it's cool!
			Relative URL: /path/to/resource
			URL without Scheme: www.example.com
			URL without Scheme and with Path: www.example.com/resources
			URL without Scheme and Query: www.example.com/search?q=openai
			URL without Scheme, IPv4: 192.168.1.1
			URL without Scheme, IPv4, and Port: 192.168.1.1:8080
			URL without Scheme, IPv6: [2001:0db8:85a3:0000:0000:8a2e:0370:7334]
			URL without Scheme, IPv6, and Port: [2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080
			`,
			[]string{
				"http://localhost",
				"http://localhost:8000/home",
				"https://www.example.com",
				"http://www.example.com:8080/page",
				"https://www.example.com/search?q=openai",
				"https://www.example.com/resources/docs/guide",
				"https://www.example.com/about#team",
				"https://user:password@example.com",
				"https://www.example.com:8080/search?q=openai#results",
				"https://www.例子.公司.cn",
				"http://192.168.1.1/setup",
				"http://192.168.1.1:8080/setup",
				"http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
				"http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080",
				"https://example.com/products?id=50&category=books",
				"https://example.com/search?name=John%20Doe&age=30",
				"https://subdomain.example.co.uk",
				"file:///C:/path/to/file.txt",
				"mailto:user@example.com",
				"ftp://user:password@ftp.example.com:21",
				"text/plain",
				"https://www.example.com/test",
				"https://www.example.com",
				"/path/to/resource",
				"www.example.com",
				"www.example.com/resources",
				"www.example.com/search?q=openai",
				"192.168.1.1",
				"192.168.1.1:8080",
				"[2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
				"[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080",
			},
		},
	}

	for _, tc := range testCases {
		got := regex.FindAllString(tc.text, -1)

		if !equalSlices(got, tc.want) {
			t.Errorf("Extracted URLs = %v, want %v", got, tc.want)
		}
	}
}

func TestURLExtractionWithHighStrictness(t *testing.T) {
	t.Parallel()

	extractor := hqgourl.NewURLExtractor(
		hqgourl.URLExtractorWithHighStrictness(),
	)

	regex := extractor.CompileRegex()

	testCases := []struct {
		text string
		want []string
	}{
		{
			`
			Localhost URL: http://localhost
			Localhost URL with Port: http://localhost:8000/home
			Standard URL: https://www.example.com
			URL with Port: http://www.example.com:8080/page
			URL with HTTPS and Query: https://www.example.com/search?q=openai
			URL with Path: https://www.example.com/resources/docs/guide
			URL with Fragment: https://www.example.com/about#team
			URL with User Info: https://user:password@example.com
			Complex URL: https://www.example.com:8080/search?q=openai#results
			International URL: https://www.例子.公司.cn
			URL with IPv4 Address: http://192.168.1.1/setup
			URL with IPv4 and Port: http://192.168.1.1:8080/setup
			URL with IPv6 Address: http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]
			URL with IPv6 and Port: http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080
			URL with Special Characters: https://example.com/products?id=50&category=books
			URL with Encoded Characters: https://example.com/search?name=John%20Doe&age=30
			URL with Multiple Subdomains: https://subdomain.example.co.uk
			File Protocol URL: file:///C:/path/to/file.txt
			Mailto Protocol URL: mailto:user@example.com
			FTP URL with Port: ftp://user:password@ftp.example.com:21
			Data URL: data:text/plain;base64,aGVsbG8=
			URL with JavaScript Protocol: javascript:alert('Hello World');
			URL with Spaces: https://www.example.com/test space
			URL in Text: Check out this link: https://www.example.com, it's cool!
			Relative URL: /path/to/resource
			URL without Scheme: www.example.com
			URL without Scheme and with Path: www.example.com/resources
			URL without Scheme and Query: www.example.com/search?q=openai
			URL without Scheme, IPv4: 192.168.1.1
			URL without Scheme, IPv4, and Port: 192.168.1.1:8080
			URL without Scheme, IPv6: [2001:0db8:85a3:0000:0000:8a2e:0370:7334]
			URL without Scheme, IPv6, and Port: [2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080
			`,
			[]string{
				"http://localhost",
				"http://localhost:8000/home",
				"https://www.example.com",
				"http://www.example.com:8080/page",
				"https://www.example.com/search?q=openai",
				"https://www.example.com/resources/docs/guide",
				"https://www.example.com/about#team",
				"https://user:password@example.com",
				"https://www.example.com:8080/search?q=openai#results",
				"https://www.例子.公司.cn",
				"http://192.168.1.1/setup",
				"http://192.168.1.1:8080/setup",
				"http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
				"http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080",
				"https://example.com/products?id=50&category=books",
				"https://example.com/search?name=John%20Doe&age=30",
				"https://subdomain.example.co.uk",
				"file:///C:/path/to/file.txt",
				"mailto:user@example.com",
				"ftp://user:password@ftp.example.com:21",
				"https://www.example.com/test",
				"https://www.example.com",
			},
		},
	}

	for _, tc := range testCases {
		got := regex.FindAllString(tc.text, -1)

		if !equalSlices(got, tc.want) {
			t.Errorf("Extracted URLs = %v, want %v", got, tc.want)
		}
	}
}

func TestURLExtractionWithHost(t *testing.T) {
	t.Parallel()

	extractor := hqgourl.NewURLExtractor(
		hqgourl.URLExtractorWithMediumStrictness(),
	)

	regex := extractor.CompileRegex()

	testCases := []struct {
		text string
		want []string
	}{
		{
			`
			Localhost URL: http://localhost
			Localhost URL with Port: http://localhost:8000/home
			Standard URL: https://www.example.com
			URL with Port: http://www.example.com:8080/page
			URL with HTTPS and Query: https://www.example.com/search?q=openai
			URL with Path: https://www.example.com/resources/docs/guide
			URL with Fragment: https://www.example.com/about#team
			URL with User Info: https://user:password@example.com
			Complex URL: https://www.example.com:8080/search?q=openai#results
			International URL: https://www.例子.公司.cn
			URL with IPv4 Address: http://192.168.1.1/setup
			URL with IPv4 and Port: http://192.168.1.1:8080/setup
			URL with IPv6 Address: http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]
			URL with IPv6 and Port: http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080
			URL with Special Characters: https://example.com/products?id=50&category=books
			URL with Encoded Characters: https://example.com/search?name=John%20Doe&age=30
			URL with Multiple Subdomains: https://subdomain.example.co.uk
			File Protocol URL: file:///C:/path/to/file.txt
			Mailto Protocol URL: mailto:user@example.com
			FTP URL with Port: ftp://user:password@ftp.example.com:21
			Data URL: data:text/plain;base64,aGVsbG8=
			URL with JavaScript Protocol: javascript:alert('Hello World');
			URL with Spaces: https://www.example.com/test space
			URL in Text: Check out this link: https://www.example.com, it's cool!
			Relative URL: /path/to/resource
			URL without Scheme: www.example.com
			URL without Scheme and with Path: www.example.com/resources
			URL without Scheme and Query: www.example.com/search?q=openai
			URL without Scheme, IPv4: 192.168.1.1
			URL without Scheme, IPv4, and Port: 192.168.1.1:8080
			URL without Scheme, IPv6: [2001:0db8:85a3:0000:0000:8a2e:0370:7334]
			URL without Scheme, IPv6, and Port: [2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080
			`,
			[]string{
				"http://localhost",
				"http://localhost:8000/home",
				"https://www.example.com",
				"http://www.example.com:8080/page",
				"https://www.example.com/search?q=openai",
				"https://www.example.com/resources/docs/guide",
				"https://www.example.com/about#team",
				"https://user:password@example.com",
				"https://www.example.com:8080/search?q=openai#results",
				"https://www.例子.公司.cn",
				"http://192.168.1.1/setup",
				"http://192.168.1.1:8080/setup",
				"http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
				"http://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080",
				"https://example.com/products?id=50&category=books",
				"https://example.com/search?name=John%20Doe&age=30",
				"https://subdomain.example.co.uk",
				"file:///C:/path/to/file.txt",
				"mailto:user@example.com",
				"ftp://user:password@ftp.example.com:21",
				"https://www.example.com/test",
				"https://www.example.com",
				"www.example.com",
				"www.example.com/resources",
				"www.example.com/search?q=openai",
				"192.168.1.1",
				"192.168.1.1:8080",
				"[2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
				"[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080",
			},
		},
	}

	for _, tc := range testCases {
		got := regex.FindAllString(tc.text, -1)

		if !equalSlices(got, tc.want) {
			t.Errorf("Extracted URLs = %v, want %v", got, tc.want)
		}
	}
}

// equalSlices checks if two slices of strings are equal.
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
