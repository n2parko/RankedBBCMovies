from django.views.generic import ListView
from models import Movie
from movies import get_movies


# MovieListView is a listview based on the Movie model
# on every request it will invoke the 'get_movies' function
# in order to refresh the list of movies available
# TODO: Add a filter based on the expiration date
class MovieListView(ListView):
    model = Movie

    def get_context_data(self, **kwargs):
        get_movies()
        return super(MovieListView, self).get_context_data()