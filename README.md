# <img src="./avatar.png" width="48"> Lunch Roulette
A slack bot to randomize lunch buddies

## What this is
- A slack bot using GO to schedule random lunches between members every period of time

## How to use
**There are 4 commands available in the app:**

help | feed
--- | ---
![help](./screenshots/help.gif) | ![feed](./screenshots/feed.gif)

stats | delete
--- | ---
![stats](./screenshots/stats.gif) | ![delete](./screenshots/delete.gif)


## How to install
### What you need before deploying / developing
- Create a new slack app in your slack workspace [link](https://api.slack.com/apps)
- Add a `.env` file following the example in the [.env_sample](.env_sample) *Note:* slack token is *Bot User OAuth Access Token* found in [https://api.slack.com/apps/<app_id>/oauth?](https://api.slack.com/apps/<app_id>/oauth?)

### Deployment
- Deploy the service wherever you want using the dockerfile included
- Make sure you have a proxy pointing to your `PORT` in [.env](env_sample) file
- Add the url pointing to the proxy to interactive messages request url field found here -> [https://api.slack.com/apps/<app_id>/interactive-messages?](https://api.slack.com/apps/<app_id>/interactive-messages?)

### Heroku
- you can also deploy using heroku by following [this guide](https://devcenter.heroku.com/articles/getting-started-with-go#deploy-the-app)
- dont forget to add the `SLACK_TOKEN` and `DB_NAME` to env vars
- **NOTE:** Heroku adds its own port to env so you should **not** do it

### Development
- make sure you have **docker version: 19.x+** installed
- run `docker-compose up --build` to launch service
- add the link found in [localhost:4040](http://localhost:4040) to [https://api.slack.com/apps/<app_id>/interactive-messages?](https://api.slack.com/apps/<app_id>/interactive-messages?) interactive messages request url field found here -> [https://api.slack.com/apps/<app_id>/interactive-messages?](https://api.slack.com/apps/<app_id>/interactive-messages?)


## Technologies used
- Golang
- [boltDB](https://github.com/boltdb/bolt)
- [shomali11/slacker](https://github.com/shomali11/slacker)
- [nlopes/slack](https://github.com/nlopes/slack)
- Docker
- Docker-compose