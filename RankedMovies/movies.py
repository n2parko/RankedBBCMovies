from django.conf import settings
from django.utils.timezone import utc
import requests
import dateutil.parser
import dateutil.tz
from datetime import datetime, timedelta


# Fetches movies from BBC iPlayer Program in JSON format
from RankedMovies.models import Movie


def get_movies():
    last_movie = Movie.objects.last()
    if last_movie and (datetime.utcnow().replace(tzinfo=utc) - last_movie.timestamp_check.replace(tzinfo=utc)) < timedelta(seconds=60):
        Movie.objects.filter(pid=last_movie.pid).update(timestamp_check=datetime.utcnow().replace(tzinfo=utc))
        return None
    response = requests.get(settings.BBC_MOVIES_URL)
    if response.status_code == 200:
        json = response.json()
        episodes = json.get('episodes', None)
        if episodes:
            for episode in episodes:
                print 'ss'
                programme = episode['programme']
                tmdb_info = get_tmdb_info(programme['title'])[0]
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


def get_tmdb_info(title):
    response = requests.get(settings.TMDB_BASE_URL.format('search', 'movie'), params={'query': title})
    if response.status_code == 200:
        json = response.json()
        return json.get('results', None)
    return None
