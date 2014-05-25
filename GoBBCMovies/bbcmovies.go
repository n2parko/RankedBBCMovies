package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"sync"
	"text/template"
	"time"
)

// Movie structure defines the needed
// information about the movies in BBC iPlayer
type Movie struct {
	Pid            string
	Title          string
	ShortSynopsis  string
	Image          string
	VoteAverage    float64
	VoteCount      float64
	AvailableUntil time.Time
}

// Movies is a structure with the list of
// episodes available and a lock for concurrent
// reading/writing
type Movies struct {
	lock     sync.RWMutex
	episodes []Movie
}

// Get is a recevier method for Movies,
// it will return the list of movies available
func (m *Movies) Get() []Movie {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.episodes
}

// Refresh keeps the list of episoded updated
func (m *Movies) Refresh() {
	response := GetRawJsonFrom("http://www.bbc.co.uk/tv/programmes/formats/films/player/episodes.json", url.Values{})
	episodes := response["episodes"].([]interface{})
	m.lock.Lock()
	defer m.lock.Unlock()
	m.episodes = make([]Movie, len(episodes))
	for key := range episodes {
		episode := episodes[key].(map[string]interface{})
		programme := episode["programme"].(map[string]interface{})
		tmdbInfo := GetTMDBInfo(programme["title"].(string))
		availableDate, _ := time.Parse(time.RFC3339, programme["available_until"].(string))
		m.episodes[key] = Movie{
			Pid:            programme["pid"].(string),
			Title:          programme["title"].(string),
			ShortSynopsis:  programme["short_synopsis"].(string),
			AvailableUntil: availableDate.UTC(),
			Image:          "http://image.tmdb.org/t/p/w185/" + tmdbInfo["poster_path"].(string),
			VoteAverage:    tmdbInfo["vote_average"].(float64),
			VoteCount:      tmdbInfo["vote_count"].(float64),
		}
	}
}

// BBCMovies is a global variable used to keep
// the results shared between multiple routines
// and avoides using a database for this simple task
var BBCMovies Movies

// GetRawJsonFrom returns a map[string]interface{} object
// making an HTTP Get request to the requestUrl given as input
// and passing the parameters to the request if necessary
func GetRawJsonFrom(requestUrl string, params url.Values) map[string]interface{} {
	client := &http.Client{}
	if params == nil {
		params = url.Values{}
	}
	r, _ := http.NewRequest("GET", requestUrl, bytes.NewBufferString(params.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))
	res, _ := client.Do(r)
	body, _ := ioutil.ReadAll(res.Body)
	var objmap map[string]interface{}
	json.Unmarshal(body, &objmap)
	return objmap
}

// GetTMDBInfo returns all the information that The Movie Database
// has stored for a given movie title
func GetTMDBInfo(title string) map[string]interface{} {
	params := url.Values{}
	params.Add("query", title)
	return GetRawJsonFrom("https://api.themoviedb.org/3/search/movie?api_key=21e4c20e8bf6156dbae8137852e39d8f",
		params)["results"].([]interface{})[0].(map[string]interface{})
}

// indexHandler handles the requests for the home directory "/"
// it renders a single page template located at "templates/movie_list.html"
// passing the list of movies available
func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/movie_list.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, BBCMovies.Get())
}

func main() {
	// Ensures to allocate a new Movies object
	// that will be shared between all routines
	BBCMovies = Movies{
		episodes: make([]Movie, 0),
	}
	// Sets the number of maxium goroutines to the 2*numberCPU + 1
	runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)

	// Starts a goroutine that keeps refreshing
	// the episodes available every hour
	go func() {
		for {
			BBCMovies.Refresh()
			time.Sleep(time.Duration(60 * time.Minute))
		}
	}()

	// Sets up the handlers and listen on port 8080
	http.HandleFunc("/", indexHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8080", nil)
}
