package oauth2

import (
  "container/list"
  "http"
  "regexp"
  "template"

  "appengine"
  "appengine/datastore"

  "fragspace/model"
)

func init() {
  http.HandleFunc("/oauth2/signup", SignUp)
}

type SignUpModel struct {
  ResponseType string
  ClientId string
  Msgs []string
}
func SignUp(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case "POST": signupPost(w, r)
    default: signupGet(w, r, make([]string, 0))
  }
}
func signupGet(w http.ResponseWriter, r *http.Request, msgs []string) {
  responseType, clientId, err := params(r)
  if err != nil {
    err.WriteTo(w)
    return
  }
  signupTemplate := template.Must(template.ParseFile("fragspace/oauth2/signup.xhtml"))
  signupModel := &SignUpModel{
    responseType,
    clientId,
    msgs,
  }
  if err := signupTemplate.Execute(w, signupModel); err != nil {
    c := appengine.NewContext(r)
    c.Errorf("%v", err)
    http.Error(w, err.String(), http.StatusInternalServerError)
  }
}
func signupPost(w http.ResponseWriter, r *http.Request) {
  _, _, err := params(r)
  //responseType, clientId, err := params(r)
  if err != nil {
    err.WriteTo(w)
    return
  }
  email, password := r.FormValue("email"), r.FormValue("password")
  emailRegexp := regexp.MustCompile(`^[a-z0-9._%\-+]+@[a-z0-9.\-]+\.[a-z]+$`)
  msgs := list.New()
  if !emailRegexp.MatchString(email) {
    msgs.PushBack("Invalid email address")
  }
  if len(password) < 6 {
    msgs.PushBack("Password is too short")
  }
  if msgsLen := msgs.Len() ; msgsLen > 0 {
    msgsSlice := make([]string, msgsLen)
    for i, el := 0, msgs.Front(); el != nil; i, el = i+1, el.Next() {
      msgsSlice[i] = el.Value.(string)
    }
    signupGet(w, r, msgsSlice)
  } else {
    user := &model.User{email}
    context := appengine.NewContext(r)
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
    w.Write([]byte("User created"))
  }
}
