package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/mchmarny/myevents/pkg/utils"
)

const (
	googleOAuthURL   = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	stateCookieName  = "authstate"
	userIDCookieName = "uid"
)

// OAuthLoginHandler handles oauth login
func OAuthLoginHandler(w http.ResponseWriter, r *http.Request) {
	uid := getCurrentUserID(r)
	if uid != "" {
		log.Printf("User ID from previous visit: %s", uid)
	}
	u := oauthConfig.AuthCodeURL(generateStateOauthCookie(w))
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

// OAuthCallbackHandler handles oauth callback
func OAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {

	oauthState, _ := r.Cookie(stateCookieName)

	// checking state of the callback
	if r.FormValue("state") != oauthState.Value {
		err := errors.New("invalid oauth state from Google")
		ErrorHandler(w, r, err, http.StatusInternalServerError)
		return
	}

	// parsing callback data
	data, err := getOAuthedUserData(r.FormValue("code"))
	if err != nil {
		log.Printf("Error while parsing user data %v", err)
		ErrorHandler(w, r, err, http.StatusInternalServerError)
		return
	}

	dataMap := make(map[string]interface{})
	json.Unmarshal(data, &dataMap)

	email := dataMap["email"]
	log.Printf("Email: %s", email)
	id := utils.MakeID(email.(string))

	//server resize image
	pic := dataMap["picture"]
	if pic != nil {
		dataMap["picture"] = utils.ServerSizeResizePlusPic(pic.(string), 200)
	}

	// set cookie for 30 days
	cookie := http.Cookie{
		Name:    userIDCookieName,
		Path:    "/",
		Value:   id,
		Expires: time.Now().Add(cookieDuration),
	}
	http.SetCookie(w, &cookie)

	// redirect on success
	if err := templates.ExecuteTemplate(w, "view", dataMap); err != nil {
		log.Printf("Error in view template: %s", err)
		ErrorHandler(w, r, err, http.StatusInternalServerError)
		return
	}

}

// OAuthLogoutHandler resets cookie and redirects to home page
func OAuthLogoutHandler(w http.ResponseWriter, r *http.Request) {

	uid := getCurrentUserID(r)
	log.Printf("User logging out: %s", uid)

	cookie := http.Cookie{
		Name:    userIDCookieName,
		Path:    "/",
		Value:   "",
		MaxAge:  -1,
		Expires: time.Now().Add(-longTimeAgo),
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther) // home
}

func generateStateOauthCookie(w http.ResponseWriter) string {

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:    stateCookieName,
		Value:   state,
		Expires: time.Now().Add(cookieDuration),
	}
	http.SetCookie(w, &cookie)

	return state
}

func getOAuthedUserData(code string) ([]byte, error) {

	// exchange code
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Got wrong exchange code: %v", err)
	}

	// user info
	response, err := http.Get(googleOAuthURL + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Error getting user info: %v", err)
	}
	defer response.Body.Close()

	// parse body
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %v", err)
	}

	return contents, nil
}
