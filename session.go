package main

import (
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/zmb3/spotify"
)

// these variables are set during the build process
// using the specified envrionment varaibles
var clientID string     //SPOTIFY_CLIENT_ID
var clientSecret string //SPOTIFY_CLIENT_SECRET
var ports = []int{47687, 47688, 47689, 47690}

// scopes is the scopes that Spotr requires to function
var scopes = []string{
	spotify.ScopePlaylistReadPrivate,
	spotify.ScopePlaylistModifyPublic,
	spotify.ScopePlaylistModifyPrivate,
	spotify.ScopePlaylistReadCollaborative,
	//spotify.ScopeUserFollowModify,
	//spotify.ScopeUserFollowRead,
	spotify.ScopeUserLibraryModify,
	spotify.ScopeUserLibraryRead,
	spotify.ScopeUserReadPrivate,
	//spotify.ScopeUserReadEmail,
}

// session manages the spotify session and authentication
type session struct {
	LastError error

	id     string
	auth   *spotify.Authenticator
	client *spotify.Client

	port       int
	cbListener net.Listener
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

	mux := http.NewServeMux()
	mux.HandleFunc("/auth-callback", s.authCallback)

	var err error
	for i, port := range ports {
		s.cbListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if nil == err {
			s.port = port
			break
		} else if i == len(ports)-1 {
			panic("cannot bind callback listener, no vailable ports")
		}
	}

	go http.Serve(s.cbListener, mux)

	auth := spotify.NewAuthenticator(s.getRedirectURL(), scopes...)
	s.auth = &auth
	s.auth.SetAuthInfo(clientID, clientSecret)

}

func (s *session) ShouldTryAgain(err error) bool {
	if nil == err {
		return false
	}
	switch e := err.(type) {
	case spotify.Error:
		if e.Status == 429 {
			fmt.Printf("waiting...  \r")
			time.Sleep(time.Second * 2)
			return true
		}

	}
	return false
}

// Client is a getter for the current session spotify client (can be nil)
func (s *session) Client() *spotify.Client {
	return s.client
}

// SearchTracks searches for track results based on the given spotify
// query, and collects up to 'pages' number of result pages (-1 for all)
func (s *session) SearchTracks(query string, pages int) []spotify.FullTrack {

	var tracks []spotify.FullTrack
	var results *spotify.SearchResult
	var err error
	limit := pages * 20

	for {
		fmt.Printf("searching...\r")
		results, err = s.Client().SearchOpt(
			query,
			spotify.SearchTypeTrack,
			&spotify.Options{
				Limit: &limit,
			},
		)
		if Session.ShouldTryAgain(err) {
			continue
		}
		break
	}
	if err != nil {
		fmt.Println("%s", err)
		return tracks
	}

	//fmt.Printf(" [%04d]\n", results.Tracks.Total)

	for i := 0; i < pages || pages == -1; i++ {
		for {
			tracks = append(tracks, results.Tracks.Tracks...)
			err = s.Client().NextTrackResults(results)
			if err == spotify.ErrNoMorePages {
				return tracks
			}
			if s.ShouldTryAgain(err) {
				continue
			}
			if nil != err {
				fmt.Printf("failed to get next result page for %s: %s", query, err)
				return tracks
			}
			break
		}
	}

	return tracks

}

// IsAuthenticated returns true if this session is logged in successfully
func (s *session) IsAuthenticated() bool {
	return (s.client != nil)
}

func (s *session) authCallback(w http.ResponseWriter, r *http.Request) {

	token, err := s.auth.Token(s.id, r)
	if err != nil {
		s.LastError = err
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		return
	}

	client := s.auth.NewClient(token)
	s.client = &client

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

	url := fmt.Sprintf("http://localhost:%d/auth-callback", s.port)
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
	//Console.Debug("auth at: " + url)

	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		fmt.Printf("please visit this URL to login: %s\n", url)
	}

	return nil

}

func (s *session) Logout() error {
	s.client = nil
	return nil
}
