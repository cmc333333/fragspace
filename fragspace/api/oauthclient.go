package api

import (
  "http"
  "regexp"
  "url"

  "appengine"
  "appengine/datastore"

  fhttp "fragspace/http"
  "fragspace/oauth2"
)

func init() {
  o := new(oauthclientHandler)
  o.Self = o
  http.Handle("/api/oauthclient", o)
}

type oauthclientHandler struct { fhttp.BaseHandler }

type oauthclientReq struct {
  Redirect string `json:"redirect_uri"`
  Name string `json:"name"`
  Email string `json:"email"`
}
type oauthclientRes struct {
  ClientId string `json:"client_id"`
  ClientSecret []byte `json:"client_secret"`
}

func (handler *oauthclientHandler) Post(r *fhttp.JsonRequest) fhttp.Response {
  post := new(oauthclientReq)
  if err := r.Extract(post); err != nil || post.Redirect == "" || post.Name == "" || post.Email == "" {
    return fhttp.UserError("invalid json")
  }
  emailRegexp := regexp.MustCompile(`^[a-z0-9._%\-+]+@[a-z0-9.\-]+\.[a-z]+$`)
  if !emailRegexp.MatchString(post.Email) {
    return fhttp.UserError("invalid email address")
  }
  if _, err := url.ParseRequest(post.Redirect); err != nil {
    return fhttp.UserError("invalid redirect uri")
  }

  context := appengine.NewContext((*http.Request)(r))
  client := oauth2.NewClient(post.Redirect, post.Name, post.Email)
  clientKey := datastore.NewKey(context, "OAuthClient", client.Id, 0, nil)
  if _, err := datastore.Put(context, clientKey, client); err != nil {
    return fhttp.ServerError(err.String())
  }
  return fhttp.JsonResponse{oauthclientRes{
    client.Id,
    client.Secret,
  }}
}
