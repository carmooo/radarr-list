# radarr-list
Connect radarr to lists from popular sources. \
Built in go using [chi](https://github.com/go-chi/chi/tree/master) and [colly](https://github.com/gocolly/colly) libraries.

## Providers
- Letterboxd

## Usage
You can use the `/letterboxd` endpoint to get a json formatted output that can be used in radarr. \
To call the endpoint you just need to add you list path from the letterboxd url after the `/letterboxd`. \
Example when running locally on port `3000`: for list `https://letterboxd.com/hepburnluv/list/classic-movies-for-beginners/` you should make a `GET` request to `localhost:3000/letterboxd/hepburnluv/list/classic-movies-for-beginners/`. \
To add it to radarr just add a list in the Steven Lu Custom format and paste the URL there.


## Run
```
make
```