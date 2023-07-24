package schemes

// SchemesUnofficial is a sorted list of some well-known url schemes which
// aren't officially registered just yet. They tend to correspond to software.
//
// Mostly collected from https://en.wikipedia.org/wiki/List_of_URI_schemes#Unofficial_but_common_URI_schemes.
var SchemesUnofficial = []string{
	`gemini`,        // gemini
	`jdbc`,          // Java database Connectivity
	`moz-extension`, // Firefox extension
	`postgres`,      // PostgreSQL (short form)
	`postgresql`,    // PostgreSQL
	`slack`,         // Slack
	`zoommtg`,       // Zoom (desktop)
	`zoomus`,        // Zoom (mobile)
}
