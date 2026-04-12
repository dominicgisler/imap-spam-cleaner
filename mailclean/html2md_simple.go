package mailclean

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// HTMLToSimpleMarkdown converts HTML into a simple Markdown-like string suitable for LLM input.
// It intentionally stays small and dependency-free beyond golang.org/x/net/html.
func HTMLToSimpleMarkdown(r io.Reader) (string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	inlineBuffer := func(s string) {
		s = strings.ReplaceAll(s, "\r\n", " ")
		s = strings.ReplaceAll(s, "\n", " ")
		s = strings.Join(strings.Fields(s), " ")
		if s != "" {
			tail := b.String()
			if b.Len() > 0 && !strings.HasSuffix(tail, "\n") && !strings.HasSuffix(tail, " ") {
				b.WriteString(" ")
			}
			b.WriteString(s)
		}
	}

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "script", "style", "head", "meta", "noscript", "iframe":
				return
			case "br":
				b.WriteString("\n")
				return
			case "p":
				if b.Len() > 0 {
					b.WriteString("\n\n")
				}
			case "h1":
				b.WriteString("\n\n# ")
			case "h2":
				b.WriteString("\n\n## ")
			case "h3":
				b.WriteString("\n\n### ")
			case "ul":
				if b.Len() > 0 {
					b.WriteString("\n")
				}
			case "ol":
				if b.Len() > 0 {
					b.WriteString("\n")
				}
			case "li":
				b.WriteString("\n- ")
			case "pre", "code":
				if b.Len() > 0 && !strings.HasSuffix(b.String(), "\n\n") {
					b.WriteString("\n\n")
				}
				b.WriteString("```\n")
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						b.WriteString(c.Data)
					} else {
						walk(c)
					}
				}
				b.WriteString("\n```\n\n")
				return
			case "a":
				href := ""
				for _, a := range n.Attr {
					if a.Key == "href" {
						href = a.Val
						break
					}
				}
				var txtBuilder strings.Builder
				var collectText func(*html.Node)
				collectText = func(nn *html.Node) {
					if nn.Type == html.TextNode {
						txtBuilder.WriteString(nn.Data)
					}
					for c := nn.FirstChild; c != nil; c = c.NextSibling {
						collectText(c)
					}
				}
				collectText(n)
				linkText := strings.TrimSpace(txtBuilder.String())
				linkText = strings.Join(strings.Fields(linkText), " ")
				needSpace := b.Len() > 0 && !strings.HasSuffix(b.String(), "\n") && !strings.HasSuffix(b.String(), " ")
				if linkText == "" && href != "" {
					if needSpace {
						b.WriteString(" ")
					}
					b.WriteString(href)
				} else if href != "" {
					if needSpace {
						b.WriteString(" ")
					}
					b.WriteString("[" + linkText + "](" + href + ")")
				} else {
					inlineBuffer(linkText)
				}
				return
			}
		}

		if n.Type == html.TextNode {
			inlineBuffer(n.Data)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}

		if n.Type == html.ElementNode {
			switch n.Data {
			case "h1", "h2", "h3", "p", "ul", "ol":
				b.WriteString("\n\n")
			}
		}
	}

	walk(doc)
	out := b.String()
	out = strings.TrimSpace(out)
	for strings.Contains(out, "\n\n\n") {
		out = strings.ReplaceAll(out, "\n\n\n", "\n\n")
	}

	// Heuristic: strip content after common quoted-reply markers if they appear
	// sufficiently deep in the body, to avoid feeding boilerplate reply chains
	// to the LLM.
	//
	// For unambiguous block-level markers (e.g. "-----Original Message-----")
	// we strip regardless of position. For generic inline patterns like "On "
	// or "From:" we require them to be at the start of a new line and past a
	// minimum offset, to avoid false positives inside regular body text.
	type sepRule struct {
		sep       string
		minOffset int
	}
	rules := []sepRule{
		{"-----original message-----", 0},
		{"----- forwarded message -----", 0},
		{"\non ", 100},
		{"\nfrom:", 100},
	}
	lower := strings.ToLower(out)
	for _, rule := range rules {
		if idx := strings.Index(lower, rule.sep); idx >= rule.minOffset {
			out = strings.TrimSpace(out[:idx])
			break
		}
	}

	return out, nil
}
