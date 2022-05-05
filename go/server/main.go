package main

import (
	"context"
	"encoding/base64"
	"flag"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/sfortson/fitness-tracker/server/migrations"
	"github.com/sfortson/fitness-tracker/server/pages"
	"github.com/sfortson/fitness-tracker/server/web/app/home"
	"github.com/sfortson/fitness-tracker/server/web/app/login"
	"github.com/sfortson/fitness-tracker/server/web/app/templates"
	abclientstate "github.com/volatiletech/authboss-clientstate"
	"github.com/volatiletech/authboss/v3"
	_ "github.com/volatiletech/authboss/v3/auth"
	"github.com/volatiletech/authboss/v3/defaults"
	_ "github.com/volatiletech/authboss/v3/logout"
	_ "github.com/volatiletech/authboss/v3/recover"
	_ "github.com/volatiletech/authboss/v3/register"
	"github.com/volatiletech/authboss/v3/remember"
)

var (
	flagDebug = flag.Bool("debug", false, "output debugging information")
	// flagDebugDB  = flag.Bool("debugdb", false, "output database on each request")
	// flagDebugCTX = flag.Bool("debugctx", false, "output specific authboss related context keys on each request")
	flagAPI = flag.Bool("api", false, "configure the app to be an api instead of an html app")
)

var (
	ab        = authboss.New()
	database  = NewMemStorer()
	schemaDec = schema.NewDecoder()

	sessionStore abclientstate.SessionStorer
	cookieStore  abclientstate.CookieStorer
)

const (
	sessionCookieName = "fitness_tracker"
)

func setupAuthboss() {
	ab.Config.Paths.RootURL = "http://localhost:8000"

	// Set up our server, session and cookie storage mechanisms.
	// These are all from this package since the burden is on the
	// implementer for these.
	ab.Config.Storage.Server = database
	ab.Config.Storage.SessionState = sessionStore
	ab.Config.Storage.CookieState = cookieStore

	// This instantiates and uses every default implementation
	// in the Config.Core area that exist in the defaults package.
	// Just a convenient helper if you don't want to do anything fancy.
	defaults.SetCore(&ab.Config, *flagAPI, false)

	// Here we initialize the bodyreader as something customized in order to accept a name
	// parameter for our user as well as the standard e-mail and password.
	//
	// We also change the validation for these fields
	// to be something less secure so that we can use test data easier.
	emailRule := defaults.Rules{
		FieldName: "email", Required: true,
		MatchError: "Must be a valid e-mail address",
		MustMatch:  regexp.MustCompile(`.*@.*\.[a-z]+`),
	}
	passwordRule := defaults.Rules{
		FieldName: "password", Required: true,
		MinLength: 4,
	}
	nameRule := defaults.Rules{
		FieldName: "name", Required: true,
		MinLength: 2,
	}

	ab.Config.Core.BodyReader = defaults.HTTPBodyReader{
		ReadJSON: *flagAPI,
		Rulesets: map[string][]defaults.Rules{
			"register":    {emailRule, passwordRule, nameRule},
			"recover_end": {passwordRule},
		},
		Confirms: map[string][]string{
			"register":    {"password", authboss.ConfirmPrefix + "password"},
			"recover_end": {"password", authboss.ConfirmPrefix + "password"},
		},
		Whitelist: map[string][]string{
			"register": {"email", "name", "password"},
		},
	}

	// Initialize authboss (instantiate modules etc.)
	if err := ab.Init(); err != nil {
		panic(err)
	}
}

func dataInjector(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := layoutData(w, &r)
		r = r.WithContext(context.WithValue(r.Context(), authboss.CTXKeyData, data))
		handler.ServeHTTP(w, r)
	})
}

// layoutData is passing pointers to pointers be able to edit the current pointer
// to the request. This is still safe as it still creates a new request and doesn't
// modify the old one, it just modifies what we're pointing to in our methods so
// we're able to skip returning an *http.Request everywhere
func layoutData(w http.ResponseWriter, r **http.Request) authboss.HTMLData {
	currentUserName := ""
	userInter, err := ab.LoadCurrentUser(r)
	if userInter != nil && err == nil {
		currentUserName = userInter.(*User).Name
	}

	return authboss.HTMLData{
		"loggedin":          userInter != nil,
		"current_user_name": currentUserName,
		// "csrf_token":        nosurf.Token(*r),
		"flash_success": authboss.FlashSuccess(w, *r),
		"flash_error":   authboss.FlashError(w, *r),
	}
}

func main() {
	token := base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(64))
	log.Println(token)
	log.Println(base64.StdEncoding.DecodeString(token))
	log.Println("Init DB...")
	// database.Open()

	log.Println("Init templates...")
	templates.InitTemplates()

	log.Println("Migrating DB...")
	migrations.Migrate()

	cookieStoreKey, _ := base64.StdEncoding.DecodeString(`NpEPi8pEjKVjLGJ6kYCS+VTCzi6BUuDzU0wrwXyf5uDPArtlofn2AG6aTMiPmN3C909rsEWMNqJqhIVPGP3Exg==`)
	sessionStoreKey, _ := base64.StdEncoding.DecodeString(`AbfYwmmt8UCwUuhd9qvfNA9UCuN1cVcKJN1ofbiky6xCyyBj20whe40rJa3Su0WOWLWcPpO1taqJdsEI/65+JA==`)
	cookieStore = abclientstate.NewCookieStorer(cookieStoreKey, nil)
	cookieStore.HTTPOnly = false
	cookieStore.Secure = false
	sessionStore = abclientstate.NewSessionStorer(sessionCookieName, sessionStoreKey, nil)
	cstore := sessionStore.Store.(*sessions.CookieStore)
	cstore.Options.HttpOnly = false
	cstore.Options.Secure = false
	cstore.MaxAge(int((30 * 24 * time.Hour) / time.Second))

	// Initialize authboss
	setupAuthboss()

	// Set up our router
	schemaDec.IgnoreUnknownKeys(true)

	// Setup router
	mux := chi.NewRouter()

	// Setup Middleware
	mux.Use(logging, ab.LoadClientStateMiddleware, remember.Middleware(ab), dataInjector)

	mux.MethodFunc("GET", "/", home.HomePage)
	mux.MethodFunc("POST", "/", home.HomePage)
	mux.MethodFunc("GET", "/registration", pages.GetRegistration)
	mux.MethodFunc("POST", "/registration", pages.SubmitRegistration)
	mux.MethodFunc("GET", "/login", login.Login)
	mux.MethodFunc("POST", "/login", login.LoginPost)

	log.Println("Listening...")
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
