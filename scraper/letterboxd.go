package scraper

import (
	"context"
	"encoding/json"
	"os"
	"regexp"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/redis/go-redis/v9"
)

const LETTERBOXD_BASE_URL string = "https://letterboxd.com"
const IMDB_REGEX_PATTERN string = "(?i)imdb.com/title/(.*?)[/$]"

type StevenLuCustomMovie struct {
	Title   string `json:"title"`
	Imdb_id string `json:"imdb_id"`
}

var ctx = context.Background()

var redisURL = os.Getenv("REDIS_URL")
var rdb = redis.NewClient(&redis.Options{
	Addr: redisURL,
})

func GetMoviesFromLetterboxd(slug string) []StevenLuCustomMovie {
	r, _ := regexp.Compile(IMDB_REGEX_PATTERN)

	c := colly.NewCollector(
		colly.AllowedDomains("letterboxd.com"),
	)

	var movies []StevenLuCustomMovie

	c.OnHTML(".col-main .poster-list .film-poster", func(e *colly.HTMLElement) {
		filmURL := e.Request.AbsoluteURL(e.Attr("data-target-link"))
		if movie, found := getCachedFilmPage(filmURL); found {
			movies = append(movies, movie)
		} else {
			e.Request.Visit(filmURL)
		}
	})

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
		cacheFilmPage(e.Request.URL.String(), movie)
	})

	c.Visit(LETTERBOXD_BASE_URL + "/" + slug)

	return movies
}

func cacheFilmPage(url string, movie StevenLuCustomMovie) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	data, _ := json.Marshal(movie)
	rdb.Set(ctx, url, data, 0)
}

func getCachedFilmPage(url string) (StevenLuCustomMovie, bool) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	data, err := rdb.Get(ctx, url).Result()
	if err != nil {
		return StevenLuCustomMovie{}, false
	}
	var movie StevenLuCustomMovie
	json.Unmarshal([]byte(data), &movie)
	return movie, true
}
