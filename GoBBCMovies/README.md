Golang RankedBBCMovies
===============

A web app that pulls the current movies on BBC iPlayer (http://www.bbc.co.uk/tv/programmes/formats/films/player/episodes.json or .xml), fetches ratings for those movies from the Movie database (http://www.themoviedb.org/) and presents it in a single page app.

The app is built with entirely in pure Go, it does not use any external library, just the plain Standard Library and it keeps an internal cache layer. It doesn't rely on SQLite. The application will refresh itself every 60 minutes.

Requirements

```
Go 1.2
```

Run

```
go run bbcmovies.go
```
