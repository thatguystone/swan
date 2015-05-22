package swan

import (
	"fmt"
	"strings"
)

func ExampleFromHTML() {
	htmlIn := `<html>
		<head>
			<title> Example Title </title>
			<meta property="og:site_name" content="Example Name"/>
		</head>
		<body>
			<p>some article body with a bunch of text in it</p>
		</body>
	</html>`

	a, err := FromHTML("http://example.com/article/1", []byte(htmlIn))
	if err != nil {
		panic(err)
	}

	if a.TopNode == nil {
		panic("no article could be extracted, " +
			"but a.Doc and a.Meta are still cleaned " +
			"and can be messed with ")
	}

	// Get the document title
	fmt.Printf("Title: %s\n", a.Meta.Title)

	// Hit any open graph tags
	fmt.Printf("Site Name: %s\n", a.Meta.OpenGraph["site_name"])

	// Print out any cleaned-up HTML that was found
	html, _ := a.TopNode.Html()
	fmt.Printf("HTML: %s\n", strings.TrimSpace(html))

	// Print out any cleaned-up text that was found
	fmt.Printf("Plain: %s\n", a.CleanedText)

	// Output: Title: Example Title
	// Site Name: Example Name
	// HTML: <p>some article body with a bunch of text in it</p>
	// Plain: some article body with a bunch of text in it
}
