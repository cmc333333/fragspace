package oauth2

import (
  "http"
  "template"

  "appengine"
)

func init() {
  http.HandleFunc("/oauth2/signup", SignUp)
}

type SignUpModel struct {
  ResponseType string
  ClientId string
}
func SignUp(w http.ResponseWriter, r *http.Request) {
  responseType, clientId, err := params(r)
  if err != nil {
    err.WriteTo(w)
    return
  }
  signupTemplate := template.Must(template.ParseFile("fragspace/oauth2/signup.xhtml"))
  signupModel := &SignUpModel{
    responseType,
    clientId,
  }
  if err := signupTemplate.Execute(w, signupModel); err != nil {
    c := appengine.NewContext(r)
    c.Errorf("%v", err)
    http.Error(w, err.String(), http.StatusInternalServerError)
  }
}
