package oauth2

import (
  "bytes"
  "encoding/base64"
  "gob"
  "http"
  "time"

  "appengine"
  "appengine/datastore"
  "appengine/memcache"

  "fragspace/encryption"
  fhttp "fragspace/http"
)

type Token struct {
  User string
  Client string
  ExpiresAt int64
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
  if grantType != "authorization_code" {
    invalidRequest("grant_type must be authorization_code").WriteTo(w)
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

  //  Make sure the clientID, secret, and redirect_uri are correct
  oauthClientKey := datastore.NewKey(context, "OAuthClient", clientId, 0, nil)
  client := new(Client)
  if err := datastore.Get(context, oauthClientKey, client); 
    err != nil || base64.StdEncoding.EncodeToString(client.Secret) != clientSecret || 
    encryption.AESDecrypt(client.Redirect, "oauthclient.redirect") != redirectUri { 
    accessDenied("invalid code").WriteTo(w)
    return
  }

  t := Token{storedCode.User, storedCode.Client, time.Seconds() + 60*60*24} //  lasts for 1 day
  var buff bytes.Buffer
  if err := gob.NewEncoder(&buff).Encode(&t); err != nil {
    panic(err)
    return
  }
  encrypted := encryption.AESByteEncrypt(buff.Bytes(), encryption.ConfigKey("oauthtoken"))
  fhttp.JsonResponse{tokenResp{encrypted}}.WriteTo(w)
}

func DecodeToken(req *http.Request) *Token {
  authentication := req.Header.Get("Authentication")
  if len(authentication) < 7 || authentication[:7] != "Bearer " {
    return nil
  }
  tokenStr, err := base64.StdEncoding.DecodeString(authentication[7:])
  if err != nil {
    return nil
  }
  decrypted := encryption.AESByteDecrypt(tokenStr, encryption.ConfigKey("oauthtoken"))
  buff := bytes.NewBuffer(decrypted)
  t := new(Token)
  if err := gob.NewDecoder(buff).Decode(t); err != nil || t.ExpiresAt < time.Seconds() {
    return nil
  }
  return t
}
