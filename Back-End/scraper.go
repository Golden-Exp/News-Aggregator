package scraper

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Article struct {
	Title     string
	URL       string
	Summary   string
	Thumbnail string
	ScrapedAt time.Time
	Source    string
	Language  string
	Region    string
}

func ScrapeDinamalar() ([]Article, error) {
	var articles []Article
	seen := make(map[string]bool)

	c := colly.NewCollector(
		colly.AllowedDomains("www.dinamalar.com"),
		colly.UserAgent("Mozilla/5.0 (compatible; NewsBot/1.0)"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*dinamalar.*",
		Parallelism: 1,
		Delay:       2 * time.Second,
	})

	// Each article block — anchor tags containing news URLs
	c.OnHTML("a[href*='/news/tamil-nadu-news/']", func(e *colly.HTMLElement) {
		url := e.Attr("href")

		// Skip comment anchors and audio anchors
		if strings.Contains(url, "#") {
			return
		}

		// Normalize to absolute URL
		if !strings.HasPrefix(url, "http") {
			url = "https://www.dinamalar.com" + url
		}

		// Deduplicate
		if seen[url] {
			return
		}
		seen[url] = true

		title := strings.TrimSpace(e.Text)

		// Skip nav links and empty titles
		if title == "" || len(title) < 5 {
			return
		}

		// Get sibling summary text (the short snippet below the title)
		summary := strings.TrimSpace(e.DOM.Next().Text())

		// Get thumbnail from nearest img
		thumbnail := e.DOM.Find("img").AttrOr("src", "")
		if thumbnail == "" {
			thumbnail, _ = e.DOM.Siblings().Find("img").Attr("src")
		}

		articles = append(articles, Article{
			Title:     title,
			URL:       url,
			Summary:   summary,
			Thumbnail: thumbnail,
			ScrapedAt: time.Now(),
			Source:    "Dinamalar",
			Language:  "ta",
			Region:    "tamil_nadu",
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Dinamalar scrape error: %s — %v", r.Request.URL, err)
	})

	err := c.Visit("https://www.dinamalar.com/news/tamil-nadu-news")
	if err != nil {
		return nil, fmt.Errorf("failed to visit Dinamalar: %w", err)
	}

	return articles, nil
}

func ScrapePuthiyathalaimurai() ([]Article, error) {
	var articles []Article
	seen := make(map[string]bool)

	c := colly.NewCollector(
		colly.AllowedDomains("www.puthiyathalaimurai.com"),
		colly.UserAgent("Mozilla/5.0 (compatible; NewsBot/1.0)"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*puthiyathalaimurai.*",
		Parallelism: 1,
		Delay:       2 * time.Second,
	})

	// Each article card — h3 inside an anchor
	c.OnHTML("a[href*='/tamilnadu/']", func(e *colly.HTMLElement) {
		url := e.Attr("href")

		// Skip nav, category, and fragment links
		if url == "/tamilnadu" || strings.HasSuffix(url, "/tamilnadu") || strings.Contains(url, "#") {
			return
		}

		// Normalize to absolute
		if !strings.HasPrefix(url, "http") {
			url = "https://www.puthiyathalaimurai.com" + url
		}

		if seen[url] {
			return
		}
		seen[url] = true

		// Tamil headline is in the h3 inside the anchor
		title := strings.TrimSpace(e.DOM.Find("h3").Text())
		if title == "" {
			// fallback: the anchor text itself
			title = strings.TrimSpace(e.Text)
		}
		if len(title) < 5 {
			return
		}

		// English slug from URL — useful as a translation hint
		parts := strings.Split(url, "/tamilnadu/")
		englishSlug := ""
		if len(parts) > 1 {
			englishSlug = strings.ReplaceAll(parts[1], "-", " ")
		}

		// Author
		author := strings.TrimSpace(e.DOM.Find("[class*='author']").Text())

		// Thumbnail — img inside the card
		thumbnail := e.DOM.Find("img").AttrOr("src", "")
		if thumbnail == "" {
			thumbnail, _ = e.DOM.Siblings().Find("img").Attr("src")
		}

		articles = append(articles, Article{
			Title:     title,
			URL:       url,
			Summary:   englishSlug, // use slug as summary hint until full article fetched
			Thumbnail: thumbnail,
			ScrapedAt: time.Now(),
			Source:    "Puthiyathalaimurai",
			Language:  "ta",
			Region:    "tamil_nadu",
		})

		_ = author // store if you extend the Article struct
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("PTT scrape error: %s — %v", r.Request.URL, err)
	})

	err := c.Visit("https://www.puthiyathalaimurai.com/tamilnadu")
	if err != nil {
		return nil, fmt.Errorf("failed to visit PTT: %w", err)
	}

	return articles, nil
}
