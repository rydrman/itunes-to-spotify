package main

import (
    "fmt"
    "net/http"
    "os/exec"
    "runtime"
    "strings"

    "github.com/zmb3/spotify"
)

// these variables are set during the build process
// using the specified envrionment varaibles
var clientID string     //SPOTIFY_CLIENT_ID
var clientSecret string //SPOTIFY_CLIENT_SECRET
var port = 47687        //4potr

// scopes is the scopes that Spotr requires to function
var scopes = []string{
    spotify.ScopePlaylistReadPrivate,
    spotify.ScopePlaylistModifyPublic,
    spotify.ScopePlaylistModifyPrivate,
    spotify.ScopePlaylistReadCollaborative,
    spotify.ScopeUserFollowModify,
    spotify.ScopeUserFollowRead,
    spotify.ScopeUserLibraryModify,
    spotify.ScopeUserLibraryRead,
    spotify.ScopeUserReadPrivate,
    //spotify.ScopeUserReadEmail,
}

// session manages the spotify session and authentication
type session struct {
    id     string
    auth   *spotify.Authenticator
    client *spotify.Client

    cbServer *http.Server
}

// Session is the singleton instance managing the spotify
// connection session for Spotr
var Session *session

func initSession() error {

    if Session != nil {
        return fmt.Errorf("session already initialized")
    }

    s := &session{
        id: RandomToken(),
    }

    Session = s
    return nil

}

// start is called once at the beginning of the program
// after the ui is initialized
func (s *session) start() {

    auth := spotify.NewAuthenticator(s.getRedirectURL(), scopes...)
    s.auth = &auth
    s.auth.SetAuthInfo(clientID, clientSecret)

    mux := http.NewServeMux()
    mux.HandleFunc("/auth-callback", s.authCallback)

    s.cbServer = &http.Server{
        Addr:    fmt.Sprintf(":%d", port),
        Handler: mux,
    }

    go s.cbServer.ListenAndServe()
    Console.Debug("session server listening on localhost" + s.cbServer.Addr)
}

// Client is a getter for the current session spotify client (can be nil)
func (s *session) Client() *spotify.Client {
    return s.client
}

func (s *session) authCallback(w http.ResponseWriter, r *http.Request) {

    token, err := s.auth.Token(s.id, r)
    if err != nil {
        Console.Error(fmt.Sprintf("Login Failed: %s", err))
        http.Error(w, "Failed to get token", http.StatusInternalServerError)
        return
    }

    client := s.auth.NewClient(token)
    s.client = &client

    name := "<UNKNOWN>"
    usr, err := s.client.CurrentUser()
    if nil == err {
        name = usr.DisplayName
    }
    Console.Logf("Login Successful! Welcome, %s", name)

    // send a self closing page
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(strings.Join(
        []string{
            "<html>",
            "<head><script type=\"text/javascript\">",
            "window.close();",
            "</script></head>",
            "<body><h1>Success!</h1></body>",
            "</html>",
        }, "")))

}

func (s *session) getRedirectURL() string {

    if s.cbServer == nil {

    }

    url := fmt.Sprintf("http://localhost:%d/auth-callback", port)
    Console.Debug("redirect to: " + url)
    return url

}

func (s *session) Authenticate() error {

    if s.client != nil {
        name := "<UNKNOWN>"
        usr, err := s.client.CurrentUser()
        if nil == err {
            name = usr.DisplayName
        }
        return fmt.Errorf(
            "Already logged in as %s", name)
    }

    url := s.auth.AuthURL(s.id)
    Console.Debug("auth at: " + url)

    switch runtime.GOOS {
    case "linux":
        exec.Command("xdg-open", url).Start()
    case "darwin":
        exec.Command("open", url).Start()
    case "windows":
        exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
    default:
        Console.Log("please visit thie URL to login: " + url)
    }

    return nil

}
