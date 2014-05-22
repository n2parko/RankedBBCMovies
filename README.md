RankedBBCMovies
===============


A web app that pulls the current movies on BBC iPlayer (http://www.bbc.co.uk/tv/programmes/formats/films/player/episodes.json or .xml), fetches ratings for those movies from the Movie database (http://www.themoviedb.org/) and presents it in a single page app.

Installation
```
pip install -r requirements.txt
python manage.py syncdb (skip this step if you plan to use the given sqlite database)
```

```
├── RankedBBCMovies
│   ├── __init__.py
│   ├── settings.py
│   ├── urls.py
│   ├── wsgi.py
├── RankedMovies
│   ├── __init__.py
│   ├── admin.py
│   ├── models.py
│   ├── movies.py
│   ├── tests.py
│   ├── views.py
├── db.sqlite3
├── manage.py
├── requirements.txt
└── templates
    ├── RankedMovies
    │   └── movie_list.html
    └── static
        └── main.css
```
