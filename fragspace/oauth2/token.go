package oauth2

import (
  "crypto/rand"
  "encoding/base64"
  "http"
  "time"

  "appengine"
  "appengine/datastore"
  "appengine/memcache"

  //fhttp "fragspace/http"
)

type token struct {
  User string
  Client string
}

func init() {
  http.HandleFunc("/oauth2/token", tokenGet)
}

type tokenResp struct {
  AccessToken []byte `json:"access_token"`
}

func tokenGet(w http.ResponseWriter, r *http.Request) {
  clientId, clientSecret := r.FormValue("client_id"), r.FormValue("client_secret")
  grantType, redirectUri, codeKey := r.FormValue("grant_type"), r.FormValue("redirect_uri"), r.FormValue("code")
  if clientId == "" || clientSecret == "" || grantType == "" || redirectUri == "" || codeKey == "" {
    invalidRequest("missing one of client_id, client_secret, grant_type, redirect_uri, or code").WriteTo(w)
    return
  }
  storedCode := new(codeStruct)
  context := appengine.NewContext(r)
  if _, err := memcache.Gob.Get(context, "oauth-code-" + codeKey, storedCode); err != nil {
    accessDenied("invalid code").WriteTo(w)
    return
  }
  memcache.Delete(context, "oauth-code-" + codeKey)

  if storedCode.ExpiresAt < time.Seconds() || storedCode.Client != clientId {
    accessDenied("invalid code").WriteTo(w)
    return
  }

  //  Make sure the clientID and secret are correct
  oauthClientKey := datastore.NewKey(context, "OAuthClient", clientId, 0, nil)
  client := new(Client)
  if err := datastore.Get(context, oauthClientKey, client); 
    err != nil || base64.StdEncoding.EncodeToString(client.Secret) != clientSecret {
    accessDenied("invalid code").WriteTo(w)
    return
  }

  buff := make([]byte, 256)
  if _, err := rand.Read(buff); err != nil {
    panic(err)
  }
  /*
  token := Token{base64.StdEncoding.EncodeToString(buff), userKey}
  tokenKey := datastore.NewKey(context, "OAuthToken", token.Id, 0, nil)
  if _, err := datastore.Put(context, tokenKey, token); err != nil {
    panic(err)
  }
  fhttp.JsonResponse{tokenResp{buff}}.WriteTo(w)
  */
}
