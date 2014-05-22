from django.views.generic import ListView
from models import Movie
from movies import get_movies


# Create your views here.
class MovieListView(ListView):
    model = Movie

    def get_context_data(self, **kwargs):
        get_movies()
        return super(MovieListView, self).get_context_data()