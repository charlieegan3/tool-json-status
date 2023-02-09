# tool-json-status

This is a tool for [toolbelt](https://github.com/charlieegan3/toolbelt) which regularly gathers information about my
public activities to keep the 'live' section on [charlieegan3.com](http://charlieegan3.com) up-to-date.

* [photos.charlieegan3.com](https://photos.charlieegan3.com)
* [music.charlieegan3.com](https://music.charlieegan3.com)
* GitHub
* Strava
* Letterboxd

Example config:

```yaml
tools:
  ... 
  json-status:
    username: xxx

    strava:
      client_secret: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
      refresh_token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
      client_id: "12345"

    play_source: https://...
    post_source: https://...

    jobs:
      refresh:
        schedule: "*/10 * * * * *"
```
