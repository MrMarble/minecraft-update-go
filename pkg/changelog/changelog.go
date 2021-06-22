package changelog

import (
	"errors"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mrmarble/minecraft-update-go/internal/emoji"
)

// Bug holds fixed bug info
type Bug struct {
	ID   string
	Link string
	Desc string
}

// Changelog holds parsed minecraft version changelog
type Changelog struct {
	URL       string
	Title     string
	New       []string
	Changes   []string
	Technical []string
	Bugs      []Bug
}

func (cl *Changelog) String() string {
	result := []string{fmt.Sprintf("<a href='%s'>%s</a> <strong>%s</strong>", emoji.Link, cl.URL, html.EscapeString(cl.Title))}

	if len(cl.New) > 0 {
		result = append(result, "")
		result = append(result, fmt.Sprintf("%s <strong>New Features</strong>", emoji.New))
		for _, n := range cl.New {
			result = append(result, fmt.Sprintf(" - %s", html.EscapeString(n)))
		}
	}
	if len(cl.Changes) > 0 {
		result = append(result, "")
		result = append(result, fmt.Sprintf("%s <strong>Changes</strong>", emoji.Change))
		for _, c := range cl.Changes {
			result = append(result, fmt.Sprintf(" - %s", html.EscapeString(c)))
		}
	}
	if len(cl.Technical) > 0 {
		result = append(result, "")
		result = append(result, fmt.Sprintf("%s <strong>Technical Changes</strong>", emoji.Technical))
		for _, t := range cl.Technical {
			result = append(result, fmt.Sprintf(" - %s", html.EscapeString(t)))
		}
	}
	if len(cl.Bugs) > 0 {
		result = append(result, "")
		result = append(result, fmt.Sprintf("%s <strong>Fixed Bugs</strong>", emoji.Bug))
		for _, b := range cl.Bugs {
			result = append(result, fmt.Sprintf(" - %s: %s", b.ID, html.EscapeString(b.Desc)))
		}
	}

	return strings.Join(result, "\n")
}

func fetch(version string) (*goquery.Document, error) {
	resp, err := http.Get(URL(version))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Changelong not found")
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// FromURL creates a Changelog from url
func FromURL(url string) (*Changelog, error) {
	doc, err := fetch(url)
	if err != nil {
		return nil, err
	}

	cl := Changelog{
		URL:     URL(url),
		Changes: []string{},
	}

	doc.Find("div.page-section--first div:nth-child(1) > h1:nth-child(1)").Each(func(i int, s *goquery.Selection) {
		cl.Title = s.Text()
	})

	doc.Find("div.page-section--first h1 + ul").Each(func(i int, s *goquery.Selection) {
		sectionTitle := s.Prev().Text()
		switch {
		case strings.HasPrefix(sectionTitle, "New"):
			s.Children().Each(func(i int, s *goquery.Selection) {
				if i < 5 {
					cl.New = append(cl.New, s.Text())
				}
			})
		case strings.HasPrefix(sectionTitle, "Changes"):
			s.Children().Each(func(i int, s *goquery.Selection) {
				if i < 5 {
					cl.Changes = append(cl.Changes, s.Text())
				}
			})
		case strings.HasPrefix(sectionTitle, "Technical"):
			s.Children().Each(func(i int, s *goquery.Selection) {
				if i < 5 {
					cl.Technical = append(cl.Technical, s.Text())
				}
			})
		case strings.HasPrefix(sectionTitle, "Fixed"):
			s.Children().Each(func(i int, s *goquery.Selection) {
				if i < 5 {
					bug := Bug{
						ID:   s.Find("a").Text(),
						Link: s.Find("a").AttrOr("href", ""),
						Desc: strings.Replace(s.Contents().Last().Text(), "\u00a0- ", "", 1),
					}
					cl.Bugs = append(cl.Bugs, bug)
				}
			})

		}

	})
	return &cl, nil
}

// URL returns a changelog url for a given minecraft version
func URL(versionURL string) string {
	return fmt.Sprintf("https://www.minecraft.net/en-us/article/minecraft-%s", versionURL)
}
