from django.conf import settings
from django.utils.timezone import utc
from datetime import datetime, timedelta
from RankedMovies.models import Movie
import requests
import dateutil.parser
import dateutil.tz


# Fetches movies from BBC iPlayer Program in JSON format
# checks the timestamp on the last movie inserted (if exists),
# every 60 minutes, based on a request, it will update the database
# with new movies fetching information from tmdb using 'get_tmdb_info' function
def get_movies():
    last_movie = Movie.objects.last()
    if last_movie and (datetime.utcnow().replace(tzinfo=utc) - last_movie.timestamp_check.replace(tzinfo=utc)) < timedelta(minutes=60):
        Movie.objects.filter(pid=last_movie.pid).update(timestamp_check=datetime.utcnow().replace(tzinfo=utc))
        return None
    response = requests.get(settings.BBC_MOVIES_URL)
    if response.status_code == 200:
        json = response.json()
        episodes = json.get('episodes', None)
        if episodes:
            for episode in episodes:
                programme = episode['programme']
                tmdb_info = get_tmdb_info(programme['title'])[0]
                # To avoid duplicates and for future usage (i.e. updating the movie title or vote counts...)
                # we perform a 'get_or_create' kind of request, this method returns a tuple containing the movie
                # object and a boolean stating if the object was already in the database or not
                Movie.objects.get_or_create(
                    pid=programme['pid'],
                    defaults={
                        'title': programme['title'],
                        'short_synopsis': programme['short_synopsis'],
                        'available_until': dateutil.parser.parse(programme['available_until']).astimezone(dateutil.tz.tzutc()),
                        'image': settings.TMDB_IMAGE_URL.format(tmdb_info['poster_path']),
                        'vote_average': tmdb_info['vote_average'],
                        'vote_count': tmdb_info['vote_count'],
                        'timestamp_check': datetime.utcnow().replace(tzinfo=utc)}
                )
    return None


# Takes a movie title as input (in a string object)
# then it requests to TMDB information about the given movie
# if the response code equals 200, then it returns the full json response
def get_tmdb_info(title):
    response = requests.get(settings.TMDB_BASE_URL.format('search', 'movie'), params={'query': title})
    if response.status_code == 200:
        json = response.json()
        return json.get('results', None)
    return None
