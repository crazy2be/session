Session Library
===============

A simple package that allows persistent server-side storage of session settings. Typical usage is just

Install:
--------

	goinstall github.com/crazy2be/session

Import:
-------

	import "github.com/crazy2be/session"

Use:
----

	s := session.Get(c, r)
	s.Get("somekey")
	s.Set("somekey", "somevalue")

Do More:
--------
Go to http://gopkgdoc.appspot.com/pkg/github.com/crazy2be/session for documentation and a function reference, always up-to-date thanks to a post-receive hook from github.