
# Deploy to Heroku 
```sh
heroku login
```

```sh
heroku git:remote -a oauth2-token-inline-hook
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

# Push changes to Heroku

```sh
git add -A . && git commit -m "Enh|Fix|Feat: {Change details}" && git push heroku master && heroku logs --tail
```
