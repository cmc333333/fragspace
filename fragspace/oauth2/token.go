package oauth2

import (
  "crypto/rand"
  "encoding/base64"
  "http"

  "appengine"
  "appengine/datastore"
  "appengine/memcache"

  fhttp "fragspace/http"
)

type Token struct {
  Id string
  User *datastore.Key
}

func init() {
  http.HandleFunc("/oauth2/token", tokenGet)
}

type tokenResp struct {
  AccessToken []byte `json:"access_token"`
}

func tokenGet(w http.ResponseWriter, r *http.Request) {
  clientId, clientSecret := r.FormValue("client_id"), r.FormValue("client_secret")
  grantType, redirectUri, code := r.FormValue("grant_type"), r.FormValue("redirect_uri"), r.FormValue("code")
  if clientId == "" || clientSecret == "" || grantType == "" || redirectUri == "" || code == "" {
    invalidRequest("missing one of client_id, client_secret, grant_type, redirect_uri, or code").WriteTo(w)
    return
  }
  var userKey *datastore.Key
  context := appengine.NewContext(r)
  if _, err := memcache.Gob.Get(context, "oauth-code-" + code, userKey); err != nil {
    accessDenied("invalid code").WriteTo(w)
    return
  }
  memcache.Delete(context, "oauth-code-" + code)

  buff := make([]byte, 256)
  if _, err := rand.Read(buff); err != nil {
    panic(err)
  }
  token := Token{base64.StdEncoding.EncodeToString(buff), userKey}
  tokenKey := datastore.NewKey(context, "OAuthToken", token.Id, 0, nil)
  if _, err := datastore.Put(context, tokenKey, token); err != nil {
    panic(err)
  }
  fhttp.JsonResponse{tokenResp{buff}}.WriteTo(w)
}
