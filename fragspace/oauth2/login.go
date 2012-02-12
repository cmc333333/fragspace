package oauth2

import (
  "http"
  "strings"
  "template"
  "url"

  "appengine"
  "appengine/datastore"

  "fragspace/model"
  "fragspace/slicelib"
)

func init() {
  http.HandleFunc("/oauth2/auth", func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
      case "POST": loginPost(w, r)
      default: loginGet(w, r)
    }
  })
}

type LoginModel struct {
  ResponseType string
  ClientId string
  Msgs []string
}

func loginGet(w http.ResponseWriter, r *http.Request) {
  responseType, client, err := params(r)
  if err != nil {
    err.WriteTo(w)
    return
  }
  loginTemplate := template.Must(template.ParseFile("fragspace/oauth2/login.xhtml"))
  loginModel := &LoginModel{
    responseType,
    client.Id,
    slicelib.Filter(strings.Split(r.FormValue("msgs"), "|"), slicelib.IsNonEmpty),
  }
  if err := loginTemplate.Execute(w, loginModel); err != nil {
    c := appengine.NewContext(r)
    c.Errorf("%v",err)
    http.Error(w, err.String(), http.StatusInternalServerError)
  }
}
func loginPost(w http.ResponseWriter, r *http.Request) {
  responseType, client, err := params(r)
  if err != nil {
    err.WriteTo(w)
    return
  }
  email, password := r.FormValue("email"), r.FormValue("password")

  emailHash := model.UserEmail(email)
  context := appengine.NewContext(r)
  query := datastore.NewQuery("User").Filter("EmailHash =", emailHash)
  found := false
  var foundKey *datastore.Key
  for row := query.KeysOnly().Run(context); ; {
    key, e := row.Next(nil)
    if e == datastore.Done {
      break
    }
    query := datastore.NewQuery("Authentication").Filter("Type =", "password").Filter("User =", key)
    for authRow := query.Run(context); ; {
      var auth model.Authentication
      _, authE := authRow.Next(&auth)
      if authE == datastore.Done {
        break
      }
      if auth.PasswordHash().EqualToPassword(password) {
        found = true
        foundKey = key
      }
    }
    if found {
      key := newCodeKey(foundKey.StringID(), client.Id, context)
      http.RedirectHandler(client.redirectUrl(key), 303).ServeHTTP(w, r)
    } else {
      http.RedirectHandler("/oauth2/auth?response_type=" + url.QueryEscape(responseType) + "&client_id=" +
        url.QueryEscape(client.Id) + "&msgs=" + url.QueryEscape("Incorrect Password"), 303).ServeHTTP(w, r)
    }
  }
}
