package scraper

import (
	"regexp"

	"github.com/gocolly/colly/v2"
)

const LETTERBOXD_BASE_URL string = "https://letterboxd.com"
const IMDB_REGEX_PATTERN string = "(?i)imdb.com/title/(.*?)[/$]"

type StevenLuCustomMovie struct {
	Title   string `json:"title"`
	Imdb_id string `json:"imdb_id"`
}

func GetMoviesFromLetterboxd(slug string) []StevenLuCustomMovie {
	r, _ := regexp.Compile(IMDB_REGEX_PATTERN)

	c := colly.NewCollector(
		colly.AllowedDomains("letterboxd.com"),
	)

	c.OnHTML(".col-main .poster-list .film-poster", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("data-target-link"))
	})

	c.OnHTML(".paginate-nextprev .next", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	var movies []StevenLuCustomMovie
	c.OnHTML("html", func(e *colly.HTMLElement) {
		movie := StevenLuCustomMovie{}
		imdbLink := e.ChildAttr("a[data-track-action=IMDb]", "href")
		if imdbLink == "" {
			return
		}
		movie.Imdb_id = r.FindStringSubmatch(imdbLink)[1]
		title := e.ChildText(".headline-1")
		if title == "" {
			return
		}
		movie.Title = title
		movies = append(movies, movie)
	})

	c.Visit(LETTERBOXD_BASE_URL + "/" + slug)

	return movies
}
