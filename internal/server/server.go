package server

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"
)

type Handler struct {
}

func RunServer() error {
	handler := &Handler{}

	server := &http.Server{
		Addr:         ":8000",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("started on port :%d\n", 8000)

	return server.ListenAndServe()
}

type Credentials struct {
	Login    string
	Password string
}

type GrantRequest struct {
	ResponseType, Scope, ClientId string
	RedirectUrl                   url.URL
}

type Token struct {
	UserId string `json:"userId"`
}

func (t Token) IsExpired() bool {
	return false
}

func token(raw string) (Token, error) {
	// TODO нормальный формат токена
	token := &Token{UserId: raw}
	return *token, nil
}

// TODO isUserAuthorized middleware
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/login":
		redirectUrl := r.URL.Query().Get("redirect_url")

		tmpl := template.Must(template.ParseFiles("web/login.html"))

		rawToken, err := r.Cookie("token")
		if err == nil {
			token, err := token(rawToken.Value)

			if err == nil && !token.IsExpired() {
				if redirectUrl == "" {
					tmpl := template.Must(template.ParseFiles("web/succeed_login.html"))
					tmpl.Execute(w, nil)
					return
				}

				http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
				return
			}
		}

		if r.Method != http.MethodPost {
			if redirectUrl == "" {
				tmpl.Execute(w, struct{ Url string }{"http://localhost:8000/login"})
				return
			}

			tmpl.Execute(w, struct{ Url string }{"http://localhost:8000/login?redirect_url=" + url.QueryEscape(redirectUrl)})
			return
		}

		credentials := Credentials{
			Login:    r.FormValue("login"),
			Password: r.FormValue("password"),
		}

		if !authenticated(credentials) {
			http.Redirect(w, r, "http://localhost:8000/login", http.StatusSeeOther)
			return
		}

		expiration := time.Now().Add(7 * 24 * time.Hour)
		cookie := http.Cookie{Name: "token", Value: credentials.Login, Expires: expiration}
		http.SetCookie(w, &cookie)

		if redirectUrl == "" {
			tmpl = template.Must(template.ParseFiles("web/succeed_login.html"))
			tmpl.Execute(w, nil)
			return
		}

		http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
	case "/oauth/authorize":
		//request example:
		//    https://authorization-server.com/authorize?
		//    response_type=code
		//    &client_id=lw5M3hczyo_eWFrEtNL2dxcN
		//    &redirect_uri=https://www.oauth.com/playground/authorization-code.html
		//    &scope=photo+offline_access
		//    &state=o6Dcw7t5QNLI-n56

		selfRedirectUrl := url.QueryEscape("http://" + r.Host + r.URL.String())
		loginUrl, _ := url.Parse("http://localhost:8000/login?redirect_url=" + selfRedirectUrl)

		rawToken, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, loginUrl.String(), http.StatusSeeOther)
			return
		}

		token, err := token(rawToken.Value)
		if err != nil || token.IsExpired() {
			http.Redirect(w, r, loginUrl.String(), http.StatusSeeOther)
			return
		}

		redirectUrl := r.URL.Query().Get("redirect_url")
		if redirectUrl == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// TODO check if redirect url is registered
		// TODO make auth code and send one
		// TODO check state
		http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
	case "/oauth/token":
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusNotImplemented)
		// request:
		// client_id=CLIENT_ID
		// client_secret=CLIENT_SECRET
		// grant_type=authorization_code
		// code=AUTHORIZATION_CODE
		// redirect_uri=CALLBACK_URL

		//_, _ = w.Write([]byte(`{"access_token":"ACCESS_TOKEN","token_type":"bearer","expires_in":2592000,"refresh_token":"REFRESH_TOKEN","scope":"read","uid":100101,"info":{"name":"Mark E. Mark","email":"mark@thefunkybunch.com"}}`))
	default:
		h.NotFound(w, r)
	}
}

func authenticated(credentials Credentials) bool {
	if credentials.Login != "s" || credentials.Password != "123" {
		return false
	}

	return true
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
