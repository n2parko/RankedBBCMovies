from django.views.generic import ListView
from models import Movie
from movies import get_movies
from datetime import datetime
from django.utils.timezone import utc


# MovieListView is a listview based on the Movie model
# on every request it will invoke the 'get_movies' function
# in order to refresh the list of movies available
class MovieListView(ListView):
    model = Movie

    # Adds a filter based on the movie availability
    queryset = Movie.objects.filter(available_until__gt=datetime.utcnow().replace(tzinfo=utc))

    def get_context_data(self, **kwargs):
        get_movies()
        return super(MovieListView, self).get_context_data()