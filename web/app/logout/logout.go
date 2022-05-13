package logout

import (
	"net/http"
	"time"

	"github.com/sfortson/fitness-tracker/internal/database"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var session database.Session
	sessionToken := c.Value
	result := database.DB.Where("session_token = ?", sessionToken).First(&session)
	if result.Error != nil {
		// If the session token is not present in session map, return an unauthorized error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// remove the users session from the session map
	database.DB.Delete(&session)

	// We need to let the client know that the cookie is expired
	// In the response, we set the session token to an empty
	// value and set its expiry as the current time
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
}
