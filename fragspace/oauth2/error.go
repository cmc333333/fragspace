package oauth2

import (
  "http"

  "libs/conf"

  fhttp "fragspace/http"
)

type Error struct {
  Error string `json:"error"`
  ErrorDescription string `json:"error_description,omitempty"`
  ErrorUri string `json:"error_uri,omitempty"`
}
func (error Error) WriteTo(w http.ResponseWriter) {
  w.WriteHeader(http.StatusBadRequest)
  fhttp.JsonResponse{error}.WriteTo(w)
}

func invalidRequest(description string) Error {
  return Error{Error: "invalid_request", ErrorDescription: description}
}
func invalidClient(description string) Error {
  return Error{Error: "invalid_client", ErrorDescription: description}
}

func params(r *http.Request) (responseType string, clientId string, retError fhttp.Response) {
  responseType, clientId = r.FormValue("response_type"), r.FormValue("client_id")
  if responseType == "" {
    retError = invalidRequest("no response type")
    return
  }
  if responseType != "code" {
    retError = invalidRequest("response_type must be code")
    return
  }
  if clientId == "" {
    retError = invalidClient("no client_id")
    return
  }
  config, err := conf.ReadConfigFile("config.ini")
  if err != nil {
    retError = fhttp.ServerError("could not read config file: " + err.String())
    return
  }
  webClientId, err := config.GetString("webclient", "client_id")
  if err != nil {
    retError = fhttp.ServerError("could not read client_id: " + err.String())
    return
  }
  if clientId != webClientId {
    retError = invalidClient("")
    return
  }
  return
}
