package link

import (
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

//This represents an 'a' html tag with href (URL)
type Link struct {
	Href string
	Text string
	Code int
}

//Parse will take in an HTML document and will return a slice of links parsed from it.
func Parse(r io.Reader) ([]Link, error) {

	doc, err := html.Parse(r)

	if err != nil {
		return nil, err
	}

	//linkNodes Takes an HTML body and returns an array of HTML elements containing the <a> tag
	nodes := linkNodes(doc)

	var links []Link

	for _, node := range nodes {
		links = append(links, buildLink(node))
	}

	return links, nil

}

//buildLink takes an html node and converts it to type link
func buildLink(n *html.Node) Link {

	var ret Link

	for _, attr := range n.Attr {

		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}

	}

	//Text function takes an html node and returns the text in it
	ret.Text = text(n)

	return ret

}

//text takes an html node and returns the text in it
func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}

	var ret string

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += text(c) + " "
	}

	return strings.Join(strings.Fields(ret), " ")

}

//Takes an HTML body and returns array of HTML elements containing the <a> tag
func linkNodes(n *html.Node) []*html.Node {

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				return []*html.Node{n}
			}
		}
	}

	var ret []*html.Node

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}

//Fixlinks reads an array of Links, appends domain name where needed, returns as new array
func Fixlinks(links []Link, url string) []Link {

	//Retrieve domain name from URL (to be used later for appending to URL paths)
	zp := regexp.MustCompile(`/`)
	var temp = zp.Split(url, -1)
	var domainname = temp[0] + "//" + temp[2]

	//Fix broken URLs then create new array with the new values
	var newlinks []Link
	for i := 0; i < len(links); i++ {
		var path = links[i].Href
		if strings.HasPrefix(path, "/") {
			path = domainname + path
		}

		if strings.HasPrefix(path, "#") {
			path = domainname + path
		}

		if !strings.HasPrefix(path, "http") && path != "" {
			path = domainname + "/" + path
		}

		var newlink Link
		newlink.Href = path
		newlink.Text = links[i].Text

		newlinks = append(newlinks, newlink)
	}

	return newlinks
}
