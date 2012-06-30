package frontend

import (
  "http"

  "appengine"
  "appengine/datastore"
  "appengine/memcache"
)

func init() {
  http.HandleFunc("/authCallback", func(w http.ResponseWriter, r *http.Request) {
    WithSession(w, r, authCallback)
  })
}

func authCallback(w http.ResponseWriter, r *http.Request, oldSession *Session) {
  //  We cheat here -- font-end has direct access to the oauth data
  code := r.FormValue("code")
  context := appengine.NewContext(r)
  var userKey datastore.Key
  _, err := memcache.Gob.Get(context, "oauth-code-" + code, &userKey)
  if err != nil {
    w.Write([]byte("Invalid code"))
  } else {
    /*
    auth := model.NewOAuth2(&userKey)
    if _, err = datastore.Put(context, datastore.NewIncompleteKey(context, "Authentication", nil), auth); err != nil {
      context.Errorf("Error saving: %v", err)
      w.Write([]byte("Error saving: " + err.String()))
      return
    }
    //  replace session cookie
    oldKey := datastore.NewKey(context, "Session", oldSession.Key, 0, nil)
    datastore.Delete(context, oldKey)

    session := NewSession(context)
    oldSession.Key = session.Key
    oldSession.Token = auth.Token
    oldSession.SetCookie(w, context)

    //  redirect
    */
  }
}
