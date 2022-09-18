
# Host on Heroku
## Deploy initial project to Heroku 

```sh
heroku login
```

```sh
heroku git:remote -a oauth2-token-inline-hook
```

```sh
heroku config:set --app=oauth2-token-inline-hook JWT_AT_AUD=Https://IDMapper-TIH-Service.com
heroku config:set --app=oauth2-token-inline-hook JWT_AT_CLIENT_ID=0oa5fmydlqntI8ExQ1d7
heroku config:set --app=oauth2-token-inline-hook JWT_AT_ISS=https://star.oktapreview.com/oauth2/aus5fqoxl0AWuk8SL1d7
heroku config:set --app=oauth2-token-inline-hook JWT_AT_REQ_SCOPE=idmapper.tihservice.execute
```

```sh
git init && git add -A .
```

```sh
git commit -m "Init"
```

```sh
git push heroku master
```

```sh
heroku logs --tail
```

## Push on-going project changes to Heroku

```sh
git add -A . && git commit -m "Enh|Fix|Feat: {Change details}" && git push heroku master && heroku logs --tail
```

[Heroku CLI Commands](https://devcenter.heroku.com/articles/heroku-cli-commands)