
# =======================
heroku login

heroku git:remote -a oauth2-token-inline-hook

git init && git add -A .

git commit -m "Init"

<!-- heroku create -->

git push heroku master

heroku logs --tail
# =======================

git add -A . && git commit -m "Enh: Restart Push" && git push heroku master && heroku logs --tail

# ====
AuthZServer : OAuth2-TIH
Application :My SPA
App Client ID: 0oa2kftuxjZmCGNB91d7

Code Verifier : abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ
Code Challenge : zwBxoIOtPkc0nS4_vIltB6DVBYCzNcN-OX1Akb-OcTs

//
https://star.oktapreview.com/oauth2/aus5dzl8k7SHsqs7w1d7/v1/authorize?client_id=0oa2kftuxjZmCGNB91d7&response_type=code&response_mode=fragment&scope=openid%20profile%20offline_access&redirect_uri=http://localhost:8080/login/callback&state=83344d15-7529-42d1-bc2c-de446bc2cd10&nonce=7ee0a4af-99d0-4372-bd65-6bd2e22872c2&code_challenge_method=S256&code_challenge=zwBxoIOtPkc0nS4_vIltB6DVBYCzNcN-OX1Akb-OcTs&b3-trace-id=xtid123987

//
https://star.oktapreview.com/oauth2/aus5e1eo246IhnMfs1d7/v1/authorize?client_id=0oa2kftuxjZmCGNB91d7&response_type=code&response_mode=fragment&scope=openid%20profile%20offline_access&redirect_uri=http://localhost:8080/login/callback&state=83344d15-7529-42d1-bc2c-de446bc2cd10&nonce=7ee0a4af-99d0-4372-bd65-6bd2e22872c2&code_challenge_method=S256&code_challenge=zwBxoIOtPkc0nS4_vIltB6DVBYCzNcN-OX1Akb-OcTs&b3-trace-id=xtid123987