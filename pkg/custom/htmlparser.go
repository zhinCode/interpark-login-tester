package custom

import (
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

func GetElementByClass(n *html.Node, tag, class string) []*html.Node {
	var result []*html.Node

	var findElements func(*html.Node)
	findElements = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tag {
			for _, attr := range n.Attr {
				if attr.Key == "class" && attr.Val == class {
					result = append(result, n)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findElements(c)
		}
	}

	findElements(n)
	return result
}

func GetTextByTag(n *html.Node, tag string) []string {
	var result []string
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tag {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					result = append(result, c.Data)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	extractText(n)
	return result
}

func DecodeBody(body []byte, charset string) (string, error) {
	if charset == "EUC-KR" {
		reader := transform.NewReader(strings.NewReader(string(body)), korean.EUCKR.NewDecoder())
		decoded, err := io.ReadAll(reader)
		if err != nil {
			return "", err
		}
		return string(decoded), nil
	}
	return string(body), nil
}
