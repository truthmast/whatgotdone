package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
)

const userKitAuthCookieName = "userkit_auth_token"

type recentEntry struct {
	Author       string `json:"author"`
	Date         string `json:"date"`
	lastModified string
	Markdown     string `json:"markdown"`
}

func (s *defaultServer) indexHandler() http.HandlerFunc {
	var templates = template.Must(
		// Use custom delimiters so Go's delimiters don't clash with Vue's.
		template.New("index.html").Delims("[[", "]]").ParseFiles(
			"./web/frontend/dist/index.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		type page struct {
			Title string
		}
		p := &page{
			Title: "What Got Done",
		}
		err := templates.ExecuteTemplate(w, "index.html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s defaultServer) submitPageHandler() http.HandlerFunc {
	var templates = template.Must(
		// Use custom delimiters so Go's delimiters don't clash with Vue's.
		template.New("index.html").Delims("[[", "]]").ParseFiles(
			"./web/frontend/dist/index.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		_, err := s.loggedInUser(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		type page struct {
			Title string
		}
		p := &page{
			Title: "What Got Done - Submit Entry",
		}
		err = templates.ExecuteTemplate(w, "index.html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *defaultServer) entriesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("Failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}

		entries, err := s.datastore.All(username)
		if err != nil {
			log.Printf("Failed to retrieve entries: %s", err)
			http.Error(w, fmt.Sprintf("Failed to retrieve entries for %s", username), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(entries); err != nil {
			panic(err)
		}
	}
}

func (s *defaultServer) recentEntriesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.datastore.Users()
		if err != nil {
			log.Printf("Failed to retrieve users: %s", err)
			return
		}

		var entries []recentEntry
		for _, username := range users {
			userEntries, err := s.datastore.All(username)
			if err != nil {
				log.Printf("Failed to retrieve entries for user %s: %s", username, err)
				return
			}
			for _, entry := range userEntries {
				// Filter low-effort posts or test posts from the recent list.
				const minimumRelevantLength = 30
				if len(entry.Markdown) < minimumRelevantLength {
					continue
				}
				entries = append(entries, recentEntry{
					Author:       username,
					Date:         entry.Date,
					lastModified: entry.LastModified,
					Markdown:     entry.Markdown,
				})
			}
		}

		sort.Slice(entries, func(i, j int) bool {
			if entries[i].Date < entries[j].Date {
				return true
			}
			if entries[i].Date > entries[j].Date {
				return false
			}
			return entries[i].lastModified < entries[j].lastModified
		})

		// Reverse the order of entries.
		for i := len(entries)/2 - 1; i >= 0; i-- {
			opp := len(entries) - 1 - i
			entries[i], entries[opp] = entries[opp], entries[i]
		}

		const maxEntries = 15
		if len(entries) > maxEntries {
			entries = entries[:maxEntries]
		}

		if err := json.NewEncoder(w).Encode(entries); err != nil {
			panic(err)
		}
	}
}

func validateEntryDate(date string) bool {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}
	const whatGotDoneEpochYear = 2019
	if t.Year() < whatGotDoneEpochYear {
		return false
	}
	if t.Weekday() != time.Friday {
		return false
	}
	if t.After(thisFriday()) {
		return false
	}
	return true
}

func thisFriday() time.Time {
	t := time.Now()
	for t.Weekday() != time.Friday {
		t = t.AddDate(0, 0, 1)
	}
	return t
}

// Catchall for when no API route matches.
func (s *defaultServer) apiRootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid API path", http.StatusBadRequest)
	}
}

func (s defaultServer) loggedInUser(r *http.Request) (string, error) {
	tokenCookie, err := r.Cookie(userKitAuthCookieName)
	if err != nil {
		return "", err
	}
	return s.authenticator.UserFromAuthToken(tokenCookie.Value)
}

func usernameFromRequestPath(r *http.Request) (string, error) {
	username := mux.Vars(r)["username"]
	// If something goes wrong in a JavaScript-based client, it will send the literal string "undefined" as the username
	// when the username variable is undefined.
	if username == "undefined" || username == "" {
		return "", errors.New("Invalid username")
	}
	return username, nil
}

func dateFromRequestPath(r *http.Request) (string, error) {
	date := mux.Vars(r)["date"]
	dateFormat := "2006-01-02"
	_, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", errors.New("Invalid date format: must be YYYY-MM-DD")
	}
	return date, nil
}
