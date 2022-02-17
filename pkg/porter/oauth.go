package porter

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/barklan/cto/pkg/postgres/models"
	"github.com/barklan/cto/pkg/security"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// TODO: randomize it
var oauthStateString = "pseudo-random"

type userData struct {
	Email         string `json:"email"`
	ID            string `json:"id"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

func initOAuth() *oauth2.Config {
	googleOauthConfig := &oauth2.Config{
		RedirectURL:  os.Getenv("OAUTH_CALLBACK_URI"),
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	return googleOauthConfig
}

// singleAuth takes email and returns jwt token.
func singleAuth(base *Base, email string) string {
	client := models.Client{}
	err := base.R.Get(&client, "select * from client where email = $1", email)
	if err != nil {
		uid4, e := uuid.NewV4()
		if e != nil {
			log.Panicln("failed to generate uuid for new client", e)
		}
		u4 := uid4.String()
		client.ID = u4
		client.Active = true
		client.Email = email

		insert := "insert into client(id, email) values ($1, $2)"
		base.R.MustExec(insert, client.ID, client.Email)
	}

	jwt, err := security.CreateJWT(base.Config, email, "")
	if err != nil {
		base.Log.Error("failed to create jwt", zap.Error(err))
		return ""
	}
	return jwt
}

func handleOAuthLogin(config *oauth2.Config, w http.ResponseWriter, r *http.Request) {
	url := config.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleOAuthCallback(base *Base, config *oauth2.Config, w http.ResponseWriter, r *http.Request) {
	content, err := getUserInfo(config, r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	user := userData{}
	if err = json.Unmarshal(content, &user); err != nil {
		log.Panicln("failed to unmarshal user data", err)
	}

	jwt := singleAuth(base, user.Email)
	http.Redirect(
		w, r,
		fmt.Sprintf("%s/guest?token=%s&name=%s&project=%s",
			base.Config.Log.ServiceHostname,
			jwt, user.Email, "",
		),
		http.StatusTemporaryRedirect,
	)
}

func getUserInfo(config *oauth2.Config, state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := config.Exchange(context.TODO(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}
