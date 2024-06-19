package forumAPI

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func Encrypt(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

func VerifyPassword(s string, M *Messages) bool {
	fmt.Println(s)
	var hasLen bool
	if len(s) >= 12 {
		hasLen = true
	}

	var hasNumber, hasUpperCase, hasLowercase, hasSpecial bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpperCase = true
		case unicode.IsLower(c):
			hasLowercase = true
		case c == '#' || c == '|':
			return false
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	M.Messages = ""
	var errorMessages []string
	if !hasNumber {
		errorMessages = append(errorMessages, "un chiffre")
	}
	if !hasUpperCase {
		errorMessages = append(errorMessages, "une majuscule")
	}
	if !hasLowercase {
		errorMessages = append(errorMessages, "une minuscule")
	}
	if !hasSpecial {
		errorMessages = append(errorMessages, "un caractère spécial")
	}
	if !hasLen {
		errorMessages = append(errorMessages, "12 caractères")
	}

	if len(errorMessages) > 0 {
		M.Messages = "Le mot de passe doit contenir au moins " + strings.Join(errorMessages[:len(errorMessages)-1], ", ")
		if len(errorMessages) > 1 {
			M.Messages += " et " + errorMessages[len(errorMessages)-1]
		} else {
			M.Messages += errorMessages[len(errorMessages)-1]
		}
	} else {
		M.Messages = ""
	}
	return hasNumber && hasUpperCase && hasLowercase && hasSpecial && hasLen
}

func EmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

// Get a cookie from a specific cookie name
func GetCookie(w http.ResponseWriter, r *http.Request, name string) int {
	cookie, err := r.Cookie(name)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusBadRequest)
			return 0
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return 0
		}
	}

	userId, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "invalid cookie value", http.StatusBadRequest)
		return 0
	}

	return userId
}

// Get a coockie from a specifique coockie name
func GetCoockieAsString(w http.ResponseWriter, r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
	}
	return cookie.Value
}

// Set user id inside a coockie
func SetCookie(w http.ResponseWriter, user User) {
	cookieId := http.Cookie{
		Name:     "userId",
		Value:    strconv.Itoa(int(user.ID)),
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookieId)
	log.Printf("Cookie set: %v", cookieId) // Log the cookie

	cookieUsername := http.Cookie{
		Name:     "username",
		Value:    user.Username,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookieUsername)
	log.Printf("Cookie set: %v", cookieUsername) // Log the cookie
}

// delete coockie
func DeleteCookies(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		cookie.MaxAge = -1
		cookie.Secure = false
		http.SetCookie(w, cookie)
	}
}

func GetUserFromURL(w http.ResponseWriter, r *http.Request) string {
	// Sépare l'URL en segments
	segments := strings.Split(r.URL.Path, "/")

	// Vérifie s'il y a suffisamment de segments dans l'URL
	if len(segments) < 2 {
		return ""
	}

	// Récupère le dernier segment de l'URL (nom d'utilisateur)
	username := segments[len(segments)-1]

	return username
}
