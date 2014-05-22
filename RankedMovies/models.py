from django.db import models


class Movie(models.Model):
    pid = models.TextField(unique=True)
    title = models.TextField()
    short_synopsis = models.TextField()
    available_until = models.DateTimeField()
    image = models.URLField()
    vote_average = models.FloatField()
    vote_count = models.IntegerField()
    timestamp_check = models.DateTimeField()