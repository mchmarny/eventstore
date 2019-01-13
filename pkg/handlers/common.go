package handlers

import (
	"time"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/mchmarny/myevents/pkg/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)


var (
	// Templates for handlers
	templates   *template.Template
	oauthConfig *oauth2.Config
	longTimeAgo  = time.Duration(3650 * 24 * time.Hour)
	cookieDuration = time.Duration(30 * 24 * time.Hour)
	knownPublisherTokens = []string{}
	allowLocalPublishers = true
)

// InitHandlers initializes OAuth package
func InitHandlers() {

	// Templates
	tmpls, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error while parsing templates: %v", err)
	}
	templates = tmpls

	// Local of hosted?
	port := utils.MustGetEnv("PORT", "8080")
	baseURL := utils.MustGetEnv("EXTERNAL_URL", fmt.Sprintf("http://localhost:%s", port))
	if baseURL == "" || !strings.HasPrefix(baseURL, "http") {
		log.Fatal("baseURL must start with HTTP or HTTPS")
	}

	// OAuth
	oauthConfig = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/auth/callback", baseURL),
		ClientID:     utils.MustGetEnv("OAUTH_CLIENT_ID", ""),
		ClientSecret: utils.MustGetEnv("OAUTH_CLIENT_SECRET", ""),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	// know publishers
	tokens := utils.MustGetEnv("KNOWN_PUBLISHER_TOKENS", "")
	knownPublisherTokens = strings.Split(tokens, ",")
	allowLocalPublishers = utils.EnvVarAsBool("ALLOW_LOCAL_PUBLISH", true)

}

func getCurrentUserID(r *http.Request) string {
	c, _ := r.Cookie(userIDCookieName)
	if c != nil {
		return c.Value
	}
	return ""
}
