Session Library
===============

Summary
-------
A simple package that allows persistent server-side storage of session settings. Typical usage is just

	s := session.Get(c, r)
	s.Get("somekey")
	s.Set("somekey", "somevalue")

Install:

	goinstall github.com/crazy2be/session

Import:

	import "github.com/crazy2be/session"

Functions
---------

	type Session struct {
		ID int64
		// contains filtered or unexported fields
	}
Session represents all information associated with a user's session.

### func Generate() (s *Session)
Creates a new session object, with the ID set to a unique number. Future versions may use a hash, but the ID will always be gaurenteed to be unique. In order to actually use a session, you should use the Get() or GetExisting() methods, they are far more useful.

### func Get(c http.ResponseWriter, r *http.Request) (s *Session)
Convieniance function to get the session object for a request based on the
value of the sessionid cookie. Creates a new session if none is found.

### func GetExisting(r *http.Request) (s *Session, e os.Error)
Same as above, but only gets a session if one exists, and does not attempt to create one.

### func Load(id int64) (s *Session, err os.Error)
Loads a session from disk with the given ID. Returns an error if the session does not exist on the server, or if the file cannot be opened.

### func NewSession() (s *Session)
Allocates a new session object and returns it.

### func (s *Session) AttachTo(c http.ResponseWriter)
Should be called AT THE START, before any html is sent.

### func (s *Session) Get(name string) (value string)
Gets a key from the map. Returns a nil string if the key is empty.

### func (s *Session) GetMap() map[string]string
For advanced purposes only, use Get() or Set() whenever possible.

### func (s *Session) Save() (err os.Error)
Forces the session to be saved to disk. Note that the sessions are saved to disk on each change currently, since there are very few changes.

### func (s *Session) Set(name, value string)
