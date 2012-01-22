package oauth2

import (
  "http"
  "template"

  "appengine"
)

func init() {
  http.HandleFunc("/oauth2/auth", Login)
}

type LoginModel struct {
  ResponseType string
  ClientId string
}

func Login(w http.ResponseWriter, r *http.Request) {
  responseType, clientId, err := params(r)
  if err != nil {
    err.WriteTo(w)
    return
  }
  loginTemplate := template.Must(template.ParseFile("fragspace/oauth2/login.xhtml"))
  loginModel := &LoginModel{
    responseType,
    clientId,
  }
  if err := loginTemplate.Execute(w, loginModel); err != nil {
    c := appengine.NewContext(r)
    c.Errorf("%v",err)
    http.Error(w, err.String(), http.StatusInternalServerError)
  }
}
