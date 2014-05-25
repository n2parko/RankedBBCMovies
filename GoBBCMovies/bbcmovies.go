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

type Greeting struct {
	ID    int64
	Hello string
}

type Movie struct {
	Pid            string
	Title          string
	ShortSynopsis  string
	Image          string
	VoteAverage    float64
	VoteCount      float64
	AvailableUntil time.Time
}

type Movies struct {
	lock     sync.RWMutex
	episodes []Movie
}

var BBCMovies Movies

func (m *Movies) Get() []Movie {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.episodes
}

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

func GetTMDBInfo(title string) map[string]interface{} {
	params := url.Values{}
	params.Add("query", title)
	return GetRawJsonFrom("https://api.themoviedb.org/3/search/movie?api_key=21e4c20e8bf6156dbae8137852e39d8f",
		params)["results"].([]interface{})[0].(map[string]interface{})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/movie_list.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, BBCMovies.Get())
}

func main() {
	BBCMovies = Movies{
		episodes: make([]Movie, 0),
	}
	runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)
	go func() {
		for {
			BBCMovies.Refresh()
			time.Sleep(time.Duration(60 * time.Minute))
		}
	}()
	http.HandleFunc("/", indexHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8080", nil)
}
