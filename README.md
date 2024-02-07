# hqgourl

[![go report card](https://goreportcard.com/badge/github.com/hueristiq/hqgourl)](https://goreportcard.com/report/github.com/hueristiq/hqgourl) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/hqgourl.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hqgourl/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/hqgourl.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hqgourl/issues?q=is:issue+is:closed) [![license](https://img.shields.io/badge/license-MIT-gray.svg?color=1E90FF)](https://github.com/hueristiq/hqgourl/blob/master/LICENSE) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-1E90FF.svg) [![contribution](https://img.shields.io/badge/contributions-welcome-1E90FF.svg)](https://github.com/hueristiq/hqgourl/blob/master/CONTRIBUTING.md)

A [Go(Golang)](http://golang.org/) package for handling URLs.

## Resources

* [Features](#features)
* [Usage](#usage)
    * [Domain Parsing](#domain-parsingn)
    * [URL Parsing](#url-parsing)
    * [URL Extraction](#url-extraction)
* [Contributing](#contributing)
* [Licensing](#licensing)
* [Credits](#credits)
    * [Contributors](#contributors)
    * [Similar Projects](#similar-projects)

## Features

* **Domain Parsing:** Break down domain names into subdomains, root domains, and TLDs.
* **URL Parsing:** Extends the standard `net/url` URLs parsing with additional fields.
* **URL Extraction:** Extracts URLs from text.

## Installation

```bash
go get -v -u github.com/hueristiq/hqgourl
```

## Usage

### Domain Parsing

```go
package main

import (
    "fmt"
    "github.com/hueristiq/hqgourl"
)

func main() {
    dp := hqgourl.NewDomainParser()

    parsedDomain := dp.Parse("subdomain.example.com")
    fmt.Printf("Subdomain: %s, Root Domain: %s, TLD: %s\n", parsedDomain.Sub, parsedDomain.Root, parsedDomain.TopLevel)
}
```

### URL Parsing

```go
package main

import (
    "fmt"
    "github.com/hueristiq/hqgourl"
)

func main() {
    up := hqgourl.NewURLParser()

    parsedURL, err := up.Parse("https://subdomain.example.com:8080/path/file.txt")
    if err != nil {
        fmt.Println("Error parsing URL:", err)
        return
    }

    fmt.Printf("Subdomain: %s\n", parsedURL.Domain.Sub)
    fmt.Printf("Root Domain: %s\n", parsedURL.Domain.Root)
    fmt.Printf("TLD: %s\n", parsedURL.Domain.TopLevel)
    fmt.Printf("Port: %d\n", parsedURL.Port)
    fmt.Printf("File Extension: %s\n", parsedURL.Extension)
}
```

Set a default scheme:

```go
up := hqgourl.NewURLParser(hqgourl.URLParserWithDefaultScheme("https"))
```

### URL Extraction

You can create a `URLExtractor` instance with default settings or customize its strictness using the provided options functions:

```go
// Default extractor with low strictness
extractor := hqgourl.NewURLExtractor()

// Custom extractor with high strictness
customExtractor := hqgourl.NewURLExtractor(hqgourl.URLExtractorWithHighStrictness())
```

Once you have an URLExtractor instance, use the CompileRegex method to compile the regex based on your strictness requirement. Then, use the compiled regex to find URLs within your text:

```go
regex := extractor.CompileRegex()

text := "Visit our website at https://example.com for more information."
urls := regex.FindAllString(text, -1)

for _, url := range urls {
    fmt.Println(url)
}
```

## Contributing

[Issues](https://github.com/hueristiq/hqgourl/issues) and [Pull Requests](https://github.com/hueristiq/hqgourl/pulls) are welcome! **Check out the [contribution guidelines](https://github.com/hueristiq/hqgourl/blob/master/CONTRIBUTING.md).**

## Licensing

This utility is distributed under the [MIT license](https://github.com/hueristiq/hqgourl/blob/master/LICENSE).

## Credits

### Contributors

Thanks to the amazing [contributors](https://github.com/hueristiq/hqgourl/graphs/contributors) for keeping this project alive.

[![contributors](https://contrib.rocks/image?repo=hueristiq/hqgourl&max=500)](https://github.com/hueristiq/hqgourl/graphs/contributors)

### Similar Projects

Thanks to similar open source projects - check them out, may fit in your needs.

[DomainParser](https://github.com/Cgboal/DomainParser) ◇ [urlx](https://github.com/goware/urlx) ◇ [xurls](https://github.com/mvdan/xurls) ◇ [goware's tldomains](https://github.com/goware/tldomains) ◇ [jakewarren's tldomains](https://github.com/jakewarren/tldomains)