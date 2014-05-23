from django.db import models


# Movie is the model that will interface
# with the database layer using the given
# attributes as columns
class Movie(models.Model):
    pid = models.TextField(unique=True, db_index=True)
    title = models.TextField()
    short_synopsis = models.TextField()
    available_until = models.DateTimeField()
    image = models.URLField()
    vote_average = models.FloatField()
    vote_count = models.IntegerField()
    timestamp_check = models.DateTimeField(db_index=True)