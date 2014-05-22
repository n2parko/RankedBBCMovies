from django.conf.urls import patterns, include, url
from RankedMovies.views import MovieListView

urlpatterns = patterns('',
    url(r'^$', MovieListView.as_view(), name='home'),
)
