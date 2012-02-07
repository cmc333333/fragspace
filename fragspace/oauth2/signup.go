package oauth2

import (
  "http"
  "regexp"
  "strings"
  "template"
  "url"

  "appengine"
  "appengine/datastore"

  "fragspace/model"
  "fragspace/slicelib"
)

func init() {
  http.HandleFunc("/oauth2/signup", func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
      case "POST": signupPost(w, r)
      default: signupGet(w, r)
    }
  })
}

type SignUpModel struct {
  ResponseType string
  ClientId string
  Msgs []string
}
func signupGet(w http.ResponseWriter, r *http.Request) {
  responseType, clientId, err := params(r)
  if err != nil {
    err.WriteTo(w)
    return
  }
  signupTemplate := template.Must(template.ParseFile("fragspace/oauth2/signup.xhtml"))
  signupModel := &SignUpModel{
    responseType,
    clientId,
    slicelib.Filter(strings.Split(r.FormValue("msgs"), "|"), slicelib.IsNonEmpty),
  }
  if err := signupTemplate.Execute(w, signupModel); err != nil {
    c := appengine.NewContext(r)
    c.Errorf("%v", err)
    http.Error(w, err.String(), http.StatusInternalServerError)
  }
}
func signupPost(w http.ResponseWriter, r *http.Request) {
  responseType, clientId, err := params(r)
  if err != nil {
    err.WriteTo(w)
    return
  }
  email, password := r.FormValue("email"), r.FormValue("password")
  emailRegexp := regexp.MustCompile(`^[a-z0-9._%\-+]+@[a-z0-9.\-]+\.[a-z]+$`)
  msgs := make([]string, 0, 5)
  if !emailRegexp.MatchString(email) {
    msgs = append(msgs, "Invalid email address")
  }
  if len(password) < 6 {
    msgs = append(msgs, "Password is too short")
  }
  //  Also check if email already exists
  user := model.NewUser(email)
  context := appengine.NewContext(r)
  countExists, e := datastore.NewQuery("User").Filter("EmailHash =", user.EmailHash).Count(context)
  if e != nil {
    context.Errorf("%v", e)
    http.Error(w, e.String(), http.StatusInternalServerError)
    return
  }
  if countExists > 0 {
    msgs = append(msgs, "Email already exists")
  }

  if msgsLen := len(msgs) ; msgsLen > 0 {
    http.RedirectHandler("/oauth2/signup?response_type=" + url.QueryEscape(responseType) + "&client_id=" +
      url.QueryEscape(clientId) + "&msgs=" + url.QueryEscape(strings.Join(msgs, "|")), 303).ServeHTTP(w,r)
  } else {
    userKey, err := datastore.Put(context, datastore.NewIncompleteKey(context, "User", nil), user)
    if err != nil {
      context.Errorf("Error saving: %v", err)
      w.Write([]byte("Error saving: " + err.String()))
      return
    }
    auth := model.NewPasswordAuth(userKey, password)
    if _, err = datastore.Put(context, datastore.NewIncompleteKey(context, "Authentication", nil), auth); err != nil {
      context.Errorf("Error saving: %v", err)
      w.Write([]byte("Error saving: " + err.String()))
      return
    }
    key := newCodeKey(userKey.StringID(), clientId, context)
    http.RedirectHandler("/authCallback?code=" + url.QueryEscape(key), 303).ServeHTTP(w, r)
  }
}
