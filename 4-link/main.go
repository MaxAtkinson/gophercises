package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func (l *Link) String() string {
	return "Link: " + l.Href + ", Text: " + l.Text
}

func newLinkFromNode(n *html.Node) Link {
	var href string

	for _, attr := range n.Attr {
		if attr.Key == "href" {
			href = attr.Val
			break
		}
	}

	sb := &strings.Builder{}
	traverseText(n, sb)

	return Link{Href: href, Text: sb.String()}
}

func trimString(s string) string {
	return strings.TrimSpace(strings.Trim(s, "\n"))
}

func traverseText(n *html.Node, sb *strings.Builder) {
	if n.Type == html.TextNode {
		if len(sb.String()) > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(trimString(n.Data))
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverseText(c, sb)
	}
}

func traverseLinks(n *html.Node, allLinks []Link) []Link {
	if n.Type == html.ElementNode && n.Data == "a" {
		allLinks = append(allLinks, newLinkFromNode(n))
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		allLinks = traverseLinks(c, allLinks)
	}
	return allLinks
}

func main() {
	filesToOpen := []string{
		"ex1.html",
		"ex2.html",
		"ex3.html",
		"ex4.html",
	}
	var allLinks []Link

	for _, filename := range filesToOpen {
		file, err := os.Open(filename)

		if err != nil {
			panic("Couldn't read file!")
		}

		markup, err := html.Parse(file)

		if err != nil {
			panic("Invalid HTML in file " + filename)
		}

		links := traverseLinks(markup, allLinks)
		allLinks = append(allLinks, links...)

		fmt.Println("\n" + filename + ":")
		for idx, val := range allLinks {
			fmt.Println(fmt.Sprint(idx+1)+". ", val.String())
		}
	}
}
