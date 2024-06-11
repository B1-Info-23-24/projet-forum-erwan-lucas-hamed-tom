package forum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	GithubClientID       = "Ov23liWXnMmbsSDNj9vg"
	GithubClientSecret   = "1370752a10bd312f7a3923c218789c7403fa5e6c"
	FacebookClientID     = "1397173297667540"
	FacebookClientSecret = "0f7ef0cc1b4b51f45e9e051d1b4f8c7d"
	GoogleClientID       = "34775775715-2fe17vhipjtnspmmf8c2d3vr96t8nvjt.apps.googleusercontent.com"
	GoogleClientSecret   = "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"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<a href="/login/github">Login with GitHub</a><br><a href="/login/facebook">Login with Facebook</a><br><a href="/login/google">Login with Google</a>`)
}

func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", GithubClientID, "http://localhost:8080/callback/github")
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	accessToken := getGithubAccessToken(code)
	data := getGithubData(accessToken)

	loggedinHandler(w, r, data)
}

func getGithubAccessToken(code string) string {
	requestBodyMap := map[string]string{
		"client_id":     GithubClientID,
		"client_secret": GithubClientSecret,
		"code":          code,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(requestJSON))
	if err != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Request failed")
	}
	defer resp.Body.Close()

	respbody, _ := ioutil.ReadAll(resp.Body)
	var ghresp struct {
		AccessToken string `json:"access_token"`
	}
	json.Unmarshal(respbody, &ghresp)
	return ghresp.AccessToken
}

func getGithubData(accessToken string) string {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		log.Panic("API Request creation failed")
	}
	req.Header.Set("Authorization", "token "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Request failed")
	}
	defer resp.Body.Close()

	respbody, _ := ioutil.ReadAll(resp.Body)
	return string(respbody)
}

func facebookLoginHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf("https://www.facebook.com/v10.0/dialog/oauth?client_id=%s&redirect_uri=%s&scope=email", FacebookClientID, "http://localhost:8080/callback/facebook")
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func facebookCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	accessToken := getFacebookAccessToken(code)
	data := getFacebookData(accessToken)
	loggedinHandler(w, r, data)
}

func getFacebookAccessToken(code string) string {
	requestBodyMap := map[string]string{
		"client_id":     FacebookClientID,
		"client_secret": FacebookClientSecret,
		"code":          code,
		"redirect_uri":  "http://localhost:8080/callback/facebook",
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, err := http.NewRequest("POST", "https://graph.facebook.com/v10.0/oauth/access_token", bytes.NewBuffer(requestJSON))
	if err != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Request failed")
	}
	defer resp.Body.Close()

	respbody, _ := ioutil.ReadAll(resp.Body)
	var fbresp struct {
		AccessToken string `json:"access_token"`
	}
	json.Unmarshal(respbody, &fbresp)
	return fbresp.AccessToken
}

func getFacebookData(accessToken string) string {
	req, err := http.NewRequest("GET", "https://graph.facebook.com/me?fields=id,name,email&access_token="+accessToken, nil)
	if err != nil {
		log.Panic("API Request creation failed")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Request failed")
	}
	defer resp.Body.Close()

	respbody, _ := ioutil.ReadAll(resp.Body)
	return string(respbody)
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&scope=email&response_type=code", GoogleClientID, "http://localhost:8080/callback/google")
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	accessToken := getGoogleAccessToken(code)
	data, err := getGoogleData(accessToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting user data from Google: %v", err), http.StatusInternalServerError)
		return
	}
	googleLoggedinHandler(w, r, data)
}

func getGoogleData(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

func getGoogleAccessToken(code string) string {
	requestBodyMap := map[string]string{
		"client_id":     GoogleClientID,
		"client_secret": GoogleClientSecret,
		"code":          code,
		"grant_type":    "authorization_code",
		"redirect_uri":  "http://localhost:8080/callback/google",
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", bytes.NewBuffer(requestJSON))
	if err != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Request failed")
	}
	defer resp.Body.Close()

	respbody, _ := ioutil.ReadAll(resp.Body)
	var gresp struct {
		AccessToken string `json:"access_token"`
	}
	json.Unmarshal(respbody, &gresp)
	return gresp.AccessToken
}

func googleLoggedinHandler(w http.ResponseWriter, r *http.Request, data string) {
	if data == "" {
		http.Error(w, "UNAUTHORIZED!", http.StatusUnauthorized)
		return
	}

	var parsedData map[string]interface{}
	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		log.Panic("JSON parse error")
	}
	fmt.Println("test data", parsedData)

	// Extract the necessary fields
	email, _ := parsedData["email"].(string)
	id := fmt.Sprintf("%v", parsedData["id"])
	picture := fmt.Sprintf("%v", parsedData["picture"])

	// Print the data for debugging purposes
	fmt.Printf("Email: %s, ID: %s, Picture: %s\n", email, id, picture)

	// Extract username from email
	username := getUsernameFromEmail(email)

	user := User{
		Username: username,
		Email:    email,
		Password: id,
	}

	// Check if the user already exists in the database
	existingUser := getUserByEmail(email)
	if existingUser != nil {
		user = *existingUser // Connect the existing user
	} else {
		// Save the new user to the database
		result := DB.Create(&user)
		if result.Error != nil {
			http.Error(w, fmt.Sprintf("Error saving user to database: %v", result.Error), http.StatusInternalServerError)
			return
		}
	}

	SetCookie(w, user)
	log.Printf("User details: %+v\n", user)

	w.Header().Set("Content-Type", "application/json")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Function to extract username from email
func getUsernameFromEmail(email string) string {
	username := strings.Split(email, "@")[0]
	return username
}

// Function to get user by email from the database
func getUserByEmail(email string) *User {
	var user User
	result := DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil // User not found
	}
	return &user
}

func loggedinHandler(w http.ResponseWriter, r *http.Request, data string) {
	if data == "" {
		http.Error(w, "UNAUTHORIZED!", http.StatusUnauthorized)
		return
	}

	var parsedData map[string]interface{}
	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		log.Panic("JSON parse error")
	}
	fmt.Println("test data", parsedData)

	// Extract the necessary fields
	login := parsedData["login"].(string)
	id := fmt.Sprintf("%v", parsedData["node_id"])
	email := fmt.Sprintf("%v", parsedData["id"])
	if parsedData["email"] != nil {
		email = parsedData["email"].(string)
	}

	// Print the data for debugging purposes
	fmt.Printf("Login: %s, Email: %s\n", login, email)
	user := User{
		Username: login,
		Email:    email,
		Password: id,
	}
	if !checkExistingUser(w, r, id) {
		result := DB.Create(&user)
		if result.Error != nil {
			http.Error(w, fmt.Sprintf("Error saving user to database: %v", result.Error), http.StatusInternalServerError)
			return
		}
	} else {
		if err := DB.Where("password = ?", user.Password).First(&user).Error; err != nil {
			http.Error(w, `{"error": "User not found"}`, http.StatusUnauthorized)
			return
		}
	}
	SetCookie(w, user)
	log.Printf("User details: %+v\n", user)

	w.Header().Set("Content-Type", "application/json")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func checkExistingUser(w http.ResponseWriter, r *http.Request, password string) bool {
	var user User
	result := DB.Where("password = ?", password).First(&user)
	if result.Error == nil {
		return true
	}

	return false
}
