package main

// Client app
type Client struct {
	id, secret, callback string
}

// User resource owner
type User struct {
	id       int
	isActive bool
}
