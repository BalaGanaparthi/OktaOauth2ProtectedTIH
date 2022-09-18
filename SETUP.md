# Setup to protect Token Hook with Okta and OAuth2 Pub/Priv Keys

## Key Setup
1. Navigate to {`Admin Console`} -> `Workflow` -> `Key Management`
2. Click `Create new key` 
3. Enter unique `Key name` (Example : `TIH-OAuth2-Key`)
4. Click `Create Key`
5. Click `Actions` on the key that was just created
6. Click on `Copy public key`
   1. > Save the copied public key for later use (during app creation) : (Referred as TIH_PUBLIC_KEY)
 
## TIH Client Application (Acts as client app for invoking TIH service)
> This Application represents the client app for Oauth2 protected TIH (this app will be used by Okta when communicating with OAuth2 protected TIH service)
1. Navigate to {`Admin Console`} -> `Applications` -> `Applications`
2. Click `Create App Integration`
3. For `Sign-in method` select `OIDC - OpenID Connect`
4. For `Application type` select `Web Application`
5. Click `Next`
6. Enter `App integration name` (Example : `TIH-OAuth2-Client-App`)
7. For `Grant type` select `Client Credentials`
8. For `Assignments` -> `Controlled access` select `Skip group assignment for now`
9. Click `Save`
10. Being at `General` tab, click on `Edit` link on the `Client Credentials` section
11. For `Client authentication` select `Public key / Private key` 
12. At `PUBLIC KEYS` section, for `Configuration` select `Save keys in Okta` (This is the default selection)
13. Click `Add key` and paste the public key saved during `Key Setup` (TIH_PUBLIC_KEY) to the text area in the `Add a public key` dialog
14. Click `Done`
15. Click `Save`
16. Click `Save` for the additional dialog `Existing client secrets will no longer be used`
17. Being at `General` tab, copy `Client ID` of the application  
    1.  > Save the copied client ID for later use (during Token Inline Hook registration) : (Referred as TIH_CLIENT-APP_CLIENT_ID)
    2.  0oa5fmydlqntI8ExQ1d7

## Create new Authorization server for TIH service (API Management)
> This Authorization Server represents the Authorization Server protecting the Oauth2 protected TIH
> Provides machine-machine tokens
1. Navigate to {`Admin Console`} -> `Security` -> `API`
2. Being at `Authorization Servers` tab, click on `Add Authorization Server`
3. Enter values for `Name`, `Audience` and `Description`
   1. > Example values : Name = `IDMapper-TIH-Service`, Audience = `Https://IDMapper-TIH-Service.com` and Description = `Issuer of Oauth2 tokens for IDMapper-TIH-Service`
   2. > Note : Value provided for Audience will be validated at the TIH Service 
4. Click `Save`
5. Create scope for `IDMapper-TIH-Service`
   1. Being at the newly created Authorization Server, navigate to `Scopes` tab
   2. Click `Add Scope`
   3. Enter values for `Name`, `Display phrase` and `Description` and leave the other fields as default (un-selected)
      1. > Example values : Name `idmapper.tihservice.execute`, Display phrase = `Execute idmapper token service` and Description = `Allow the client to execute IDMapper token inline service`
   4. Click `Create`
6. Navigate to `Access Policies` tab
7. Click `Add Policy`
8. Enter values for `Name`, `Description` and for `Assign to` select `The following clients:` and add the TIH OAuth2 Client application (`TIH-OAuth2-Client-App`)
   1. > Example values : Name = `IDMapper TIH Service Policy` and Description = `IDMapper token inline hook service policy`
9. Click `Create Policy`
10. Being at the access policy just created, click on `Add Rule`
11. And enter the following values to create the rule
    1.  Enter the `Rule Name` (Example : `IDMapper TIH Service Rule`)
    2.  For `IF Grant type is` select only (`Client acting on behalf of itself`)  -> `Client Credentials`
    3.  For `Scopes requested` select `Assigned the app and a member of one of the following:` and enter `idmapper.tihservice.execute` into the textbox
    4.  Leave the rest as default
12. Click `Create Rule`
13. Copy Issuer url of the authorization server
    1.  Being at the authorization server, navigate to `Settings` tab
    2.  Click the `Metadata URI` link
    3.  Copy the value for `token_endpoint` from the metadata and save the value for later use (referred as ISSUER_URL)
    4.  https://star.oktapreview.com/oauth2/aus5fqoxl0AWuk8SL1d7/v1/token

## Deploy Token Inline Hook 
1. Get the token inline hook endpoint url that is Oauth2 protected with tokens issued from TIH service Authorization Server (Referred as TIH_ENDPOINT_URL)

## Register OAuth protected Token Inline Hook
1. Navigate to {`Admin Console`} -> `Workflow` -> `Inline Hooks`
2. Click on `Add Inline Hook` and select `Token` from the drop-down menue
3. Enter `Name` (Example : `Oauth2 Protected IDMapper`)
4. Enter `URL` (With the TIH_ENDPOINT_URL of the deployed TIH service)  
4. For `Authentication` select `OAuth 2.0`
5. For `Oauth2 Protected IDMapper` select `Use private key`
6. For `Client ID` enter the client id of the TIH client app : (Value stored previously as TIH_CLIENT-APP_CLIENT_ID) 
7. For `Public key` select the key created previously (Example : TIH-OAuth2-Key)
8. For `Token URL` enter the value saved as `${TIH_ENDPOINT_URL}`
9. For Scope enter the scope create previousl : `idmapper.tihservice.execute`
10. Click `Save`

## Create a Single Page Application that represents as a Customer facing application
> Represents user facing application
1. Navigate to {`Admin Console`} -> `Applications` -> `Applications`
2. Click `Create App Integration`
3. For `Sign-in method` select `OIDC - OpenID Connect`
4. For `Application type` select `Single-Page Application`
5. Click `Next`
6. Enter `App integration name` (Example : `ClientSPApp`)
7. For `Grant type` accept defaults (`Authorization Code` Selected)
8. For `Assignments` -> `Controlled access` select `Allow everyone in your organization to access`
9. For `Enable immediate access` accept default (`Enable immediate access with Federation Broker Mode` Selected)
10. Click `Save`
11. Being at `General` tab, copy `Client ID` of the application  
    1.  > Save the copied client ID for later use : (Referred as SP-APP_CLIENT_ID)
    2.  0oa5fwevhkX4mJhEE1d7

## Create an Authorization Server for Authenticating and Authorizing the end-users
> Represents the Authorization server interacting with end-users accessing user-facing apps
> Provides user-app tokens (ID/Access/refresh)
> Can use the Default authorization server
> Configured with the Access Policy Rules to invoke the Oauth2 protected token inline hook for token enrichment
1. Navigate to {`Admin Console`} -> `Security` -> `API`
2. Being at `Authorization Servers` tab, click on `Add Authorization Server`
3. Enter values for `Name`, `Audience` and `Description`
   1. > Example values : Name = `EndUser-AS`, Audience = `Https://End-User.com` and Description = `Issuer of end-user tokens for client facing applications`
4. Click `Save`
5. Navigate to `Access Policies` tab
7. Click `Add Policy`
8. Enter values for `Name`, `Description` and for `Assign to` select `The following clients:` and add the Client facing single page application (`ClientSPApp`)
   1. > Example values : Name = `End-User Policy` and Description = `Policy with rule to invoke OAuth2 protected TIH`
9. Click `Create Policy`
10. Being at the access policy just created, click on `Add Rule`
11. And enter the following values to create the rule
    1.  Enter the `Rule Name` (Example : `End-User Rule`)
    2.  For `THEN Use this inline hook` select Token inline hook registered previously (`Oauth2 Protected IDMapper`)
    3.  Leave the rest as default
12. Click `Create Rule`
13. Copy Issuer url of the authorization server
    1.  Being at the authorization server, navigate to `Settings` tab
    2.  Click the `Metadata URI` link
    3.  Copy the value for `authorization_endpoint` from the metadata and save the value for later use (referred as CLIENT_AUTHORIZE_URL)
    4.  https://star.oktapreview.com/oauth2/aus5fx72s7ZE7T5bu1d7/v1/authorize

## Standalone Testing
1. Generate the Authorization URL for the single page app

   1. > Template
   
   ```template 
    {CLIENT_AUTHORIZE_URL}?
    client_id={SP-APP_CLIENT_ID}&
    response_type=code&
    response_mode=fragment&
    scope={SCOPES}&
    redirect_uri={REDIRECT_URI}&
    state={STATE}&
    nonce={NONCE}&
    code_challenge_method=S256&
    code_challenge={CODE_CHALLENGE}
    ```
2. Sample Authorization URL
    Substitutions used
    ```substitutions
    CLIENT_AUTHORIZE_URL : https://star.oktapreview.com/oauth2/aus5fx72s7ZE7T5bu1d7/v1/authorize
    SP-APP_CLIENT_ID : 0oa5fwevhkX4mJhEE1d7
    SCOPES : openid%20profile%20offline_access
    REDIRECT_URI : http://localhost:8080/login/callback
    STATE
    NONCE
    CODE_VERIFIER : abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ
    CODE_CHALLENGE : zwBxoIOtPkc0nS4_vIltB6DVBYCzNcN-OX1Akb-OcTs 
    ```

    [Sample Authorization URL Link](https://star.oktapreview.com/oauth2/aus5fx72s7ZE7T5bu1d7/v1/authorize?client_id=0oa5fwevhkX4mJhEE1d7&response_type=code&response_mode=fragment&scope=openid%20profile%20offline_access&redirect_uri=http://localhost:8080/login/callback&state=83344d15-7529-42d1-bc2c-de446bc2cd10&nonce=7ee0a4af-99d0-4372-bd65-6bd2e22872c2&code_challenge_method=S256&code_challenge=zwBxoIOtPkc0nS4_vIltB6DVBYCzNcN-OX1Akb-OcTs)
  
3.  Execut the constructed URL in the browser 
4.  Get the `code` value from the browser after invoking the URL
5.  Execute the following curl or httpie to exchange the code for tokens
6.  CURL
    ```curl
    curl --location --request POST 'https://star.oktapreview.com/oauth2/aus5fx72s7ZE7T5bu1d7/v1/token' \
    --header 'Accept: application/json' \
    --header 'Content-Type: application/x-www-form-urlencoded' \
    --data-urlencode 'grant_type=authorization_code' \
    --data-urlencode 'client_id=0oa5fwevhkX4mJhEE1d7' \
    --data-urlencode 'redirect_uri=http://localhost:8080/login/callback' \
    --data-urlencode 'code=mJbnpnHRQMbDa4PhbHuVP6BfK7G4O4aMO6hZC4RCovE' \
    --data-urlencode 'code_verifier=abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ'
    ```
7.  HTTPIe
    ```http
    POST /oauth2/aus5fx72s7ZE7T5bu1d7/v1/token HTTP/1.1
    Host: star.oktapreview.com
    Accept: application/json
    Content-Type: application/x-www-form-urlencoded
    grant_type=authorization_code&client_id=0oa5fwevhkX4mJhEE1d7&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Flogin%2Fcallback&code=Oz06noF9823Y2tJpgUnKO1xDMcFgwpH2Y4WrXgTCq-s&code_verifier=abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ
    ```
8. Check the jwt tokens at jwt.io and verify if the token is enriched

9. Sample enriched Access Token (Add new claim with name `aud1` and the value with original audience value with a 1 suffix)
    ```json
    {
    "ver": 1,
    "jti": "AT.XsIhBnzWGq5rQggz1VAe1ytYYb1Y57mMDgQVMKvuEsc",
    "iss": "https://star.oktapreview.com/oauth2/aus5fx72s7ZE7T5bu1d7",
    "aud": "Https://End-User.com",
    "iat": 1663449854,
    "exp": 1663453454,
    "cid": "0oa5fwevhkX4mJhEE1d7",
    "uid": "00u1sps8sahxIfoYw1d7",
    "scp": [
        "profile",
        "openid"
    ],
    "auth_time": 1663443209,
    "sub": "bala.ganaparthi@okta.com",
    "aud1": "Https://End-User.com1"
    }
    ```

