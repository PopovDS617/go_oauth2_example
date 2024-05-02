package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golangoauth2example/internal/model"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAUTHConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8000/auth/google/callback/",
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func OAUTHGoogleLogin(w http.ResponseWriter, r *http.Request) {

	// Создать куку
	oauthState := generateStateOauthCookie(w)

	// AuthCodeURL получает состояние, которое является токеном для защиты от CSRF.
	// Должна быть указана непустая строка и должен соответствовать параметру запроса состояния в обратном вызове
	u := googleOAUTHConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func OAUTHGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Прочитать стейт из куки
	oauthState, err := r.Cookie("oauthstate")

	if err != nil {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	user, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Сохранить полученного пользователя в куке в виде jwt токена, чтобы потом прочитать его на другой странице.
	// Обычно пользователь должен быть сохранен в db и получен уже оттуда
	token, err := CreateJWTTokenFromUserData(user)

	if err != nil {
		log.Println(err)
	}

	cookie := http.Cookie{Name: "gotestsession", Value: token, Expires: time.Now().Add(20 * time.Minute), HttpOnly: true, Path: "/"}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/account/", http.StatusTemporaryRedirect)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration, Path: "/"}
	http.SetCookie(w, &cookie)

	return state
}

func getUserDataFromGoogle(code string) (model.User, error) {
	user := model.User{}

	token, err := googleOAUTHConfig.Exchange(context.Background(), code)
	if err != nil {
		return user, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return user, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&user)

	if err != nil {
		return user, fmt.Errorf("failed read response: %s", err.Error())
	}

	return user, nil
}
