package scraper

import (
	"context"
	"encoding/json"
	"log"
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
	log.Printf("Starting to scrape movies from Letterboxd for slug: %s", slug)
	r, _ := regexp.Compile(IMDB_REGEX_PATTERN)

	c := colly.NewCollector(
		colly.AllowedDomains("letterboxd.com"),
	)

	var movies []StevenLuCustomMovie

	c.OnHTML(".col-main .poster-list .film-poster", func(e *colly.HTMLElement) {
		filmURL := e.Request.AbsoluteURL(e.Attr("data-target-link"))
		log.Printf("Found film poster with URL: %s", filmURL)
		if movie, found := getCachedFilmPage(filmURL); found {
			movies = append(movies, movie)
		} else {
			e.Request.Visit(filmURL)
		}
	})

	c.OnHTML(".paginate-nextprev .next", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		movie := StevenLuCustomMovie{}
		imdbLink := e.ChildAttr("a[data-track-action=IMDb]", "href")
		if imdbLink == "" {
			log.Println("No IMDb link found, skipping")
			return
		}
		movie.Imdb_id = r.FindStringSubmatch(imdbLink)[1]
		title := e.ChildText(".headline-1")
		if title == "" {
			log.Println("No title found, skipping")
			return
		}
		movie.Title = title
		log.Printf("Scraped movie: %s (IMDb ID: %s)", movie.Title, movie.Imdb_id)
		movies = append(movies, movie)
		cacheFilmPage(e.Request.URL.String(), movie)
	})

	log.Printf("Visiting URL: %s/%s", LETTERBOXD_BASE_URL, slug)
	c.Visit(LETTERBOXD_BASE_URL + "/" + slug)

	log.Printf("Finished scraping movies for slug: %s", slug)
	return movies
}

func cacheFilmPage(url string, movie StevenLuCustomMovie) {
	log.Printf("Caching movie: %s (URL: %s)", movie.Title, url)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	data, _ := json.Marshal(movie)
	rdb.Set(ctx, url, data, 0)
}

func getCachedFilmPage(url string) (StevenLuCustomMovie, bool) {
	log.Printf("Checking cache for URL: %s", url)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	data, err := rdb.Get(ctx, url).Result()
	if err != nil {
		log.Printf("Cache miss for URL: %s", url)
		return StevenLuCustomMovie{}, false
	}
	var movie StevenLuCustomMovie
	json.Unmarshal([]byte(data), &movie)
	log.Printf("Cache hit for URL: %s", url)
	return movie, true
}
