RankedBBCMovies
===============

A web app that pulls the current movies on BBC iPlayer (http://www.bbc.co.uk/tv/programmes/formats/films/player/episodes.json or .xml), fetches ratings for those movies from the Movie database (http://www.themoviedb.org/) and presents it in a single page app.

The app is built with Django, uses SQLite as database and Twitter Bootstrap at CSS layer.

Installation
```
pip install -r requirements.txt
python manage.py syncdb (skip this step if you plan to use the given sqlite database)
```

Starting the server
```
python manage.py runserver
```
