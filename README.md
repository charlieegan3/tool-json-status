# json-charlieegan3

A rake task to keep the live content on my homepage at
[charlieegan3.com](http://charlieegan3.com) up-to-date.

Fetches data from:

* Instagram
* Twitter
* last.fm
* GitHub
* Strava
* Letterboxd

The task is currently hosted on Kubernetes as a cronjob and updates this
[status file](https://charlieegan3.com/status.json).
