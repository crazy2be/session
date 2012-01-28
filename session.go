// A simple package that allows persistent server-side storage of session settings. Typical usage is just 
//	s := session.Get(c, r)
//	s.Get("somekey")
//	s.Set("somekey", "somevalue")
package session

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	// Local Imports
	"github.com/crazy2be/httputil"
	"github.com/crazy2be/ini"
)

// The lastID given out. Since all IDs are assigned in numeric order, this should ensure there are no collisions.
var lastID int64

// Session represents all information associated with a user's session.
type Session struct {
	ID       int64
	settings map[string]string
	updated  bool
}

// Convieniance function to get the session object for a request based on the
// value of the sessionid cookie. Creates a new session if none is found.
func Get(c http.ResponseWriter, r *http.Request) (s *Session) {
	s, e := GetExisting(r)
	if e == nil {
		return
	}
	s = Generate()
	s.AttachTo(c)

	return
}

// Same as above, but only gets a session if one exists, and does not attempt to create one.
func GetExisting(r *http.Request) (s *Session, e error) {
	cookie := httputil.FindCookie(r, "sessionid")

	if cookie == nil {
		e = errors.New("No sessionid found!")
		fmt.Println(e)
		return
	}

	sid := cookie.Value
	id, e := strconv.ParseInt(sid, 10, 64)
	s, e = Load(id)
	return
}

// Allocates a new session object and returns it.
func NewSession() (s *Session) {
	s = new(Session)
	s.settings = make(map[string]string, 10)
	s.updated = false
	return s
}

// Creates a new session object, with the ID set to a unique number. Future versions may use a hash, but the ID will always be gaurenteed to be unique. In order to actually use a session, you should use the Get() or GetExisting() methods, they are far more useful.
func Generate() (s *Session) {
	s = NewSession()
	// TODO: Generate some sort of hash for the ID, rather than an int. The int would theoretically be fairly easy to guess.
	idseed := time.Now()
	// Prevent two requests during the same nanosecond from getting duplicate
	// sessionids.
	if idseed == lastID {
		idseed++
	}
	lastID = idseed

	s.ID = idseed

	return
}

// Sets a key in the map, then saves the session file.
func (s *Session) Set(name, value string) {
	s.settings[name] = value
	s.updated = true
	// Note that this will cause lag if called a lot.
	s.Save()
}

// Gets a key from the map. Returns a nil string if the key is empty.
func (s *Session) Get(name string) (value string) {
	return s.settings[name]
}

// For advanced purposes only, use Get() or Set() whenever possible.
func (s *Session) GetMap() map[string]string {
	return s.settings
}

// Should be called AT THE START, before any html is sent.
func (s *Session) AttachTo(c http.ResponseWriter) {
	// TODO: Should eventually be setting an expiration date on this...
	header := c.Header()
	header["Set-Cookie"] = append(header["Set-Cookie"], "Sessionid="+strconv.FormatInt(s.ID, 10)+"; path=/")
}

// Loads a session from disk with the given ID. Returns an error if the session does not exist on the server, or if the file cannot be opened.
func Load(id int64) (s *Session, err error) {
	s = new(Session)
	filename := "data/shared/sessions/" + strconv.FormatInt(id, 10)

	s.settings, err = ini.Load(filename)
	if err != nil {
		return
	}

	return
}

// Forces the session to be saved to disk. Note that the sessions are saved to disk on each change currently, since there are very few changes.
func (s *Session) Save() (err error) {
	filename := "data/shared/sessions/" + strconv.FormatInt(s.ID, 10)
	fmt.Println(filename)
	err = ini.Save(filename, s.settings)
	if err != nil {
		return
	}
	return
}

// Make required directories
func init() {
	os.MkdirAll("data/shared/sessions", 0755)
}
