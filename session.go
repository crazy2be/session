// A simple package that allows persistent server-side storage of session settings. Typical usage is just 
//	s := session.Get(c, r)
//	s.Get("somekey")
//	s.Set("somekey", "somevalue")
package session

import (
	"strconv"
	"http"
	"time"
	"fmt"
	"os"
	// Local Imports
	"github.com/crazy2be/ini"
	"github.com/crazy2be/httputil"
)

// The lastID given out. Since all IDs are assigned in numeric order, this should ensure there are no collisions.
var lastID int64

// Session represents all information associated with a user's session.
type Session struct {
	ID int64
	settings map[string]string
	updated bool
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
func GetExisting(r *http.Request) (s *Session, e os.Error) {
	cookie := httputil.FindCookie(r, "sessionid")
	
	if cookie == nil {
		e = os.NewError("No sessionid found!")
		fmt.Println(e)
		return
	}
	
	sid := cookie.Value
	id, e := strconv.Atoi64(sid)
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
	idseed := time.Nanoseconds()
	// Prevent two requests during the same nanosecond from getting duplicate
	// sessionids.
	if idseed == lastID {
		idseed++
	}
	lastID = idseed
	
	s.ID = idseed
	
	return
}

// Loads a session from disk with the given ID. Returns an error if the session does not exist on the server, or if the file cannot be opened.
func Load(id int64) (s* Session, err os.Error) {
	filename := "data/shared/sessions/" + strconv.Itoa64(s.ID)
	
	s.settings, err = ini.Load(filename)
	if err != nil {
		return
	}
	
	return
}

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
	header["Set-Cookie"] = append(header["Set-Cookie"], "Sessionid="+strconv.Itoa64(s.ID)+"; path=/")
}

// Forces the session to be saved to disk. Note that the sessions are saved to disk on each change currently, since there are very few changes.
func (s* Session) Save() (err os.Error) {
	filename := "data/shared/sessions/" + strconv.Itoa64(s.ID)
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