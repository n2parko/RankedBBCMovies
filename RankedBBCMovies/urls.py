from django.conf.urls import patterns, url
from RankedMovies.views import MovieListView

# urlpatterns defines all urls in this Django application
urlpatterns = patterns('',
    url(r'^$', MovieListView.as_view(), name='home'),
)
