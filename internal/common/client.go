package common

// Client app
type Client struct {
	id, secret, callbackUrl string
}

func New(id, secret, callbackUrl string) *Client {
	return &Client{
		id:          id,
		secret:      secret,
		callbackUrl: callbackUrl,
	}
}

func IsRegistered(id string) bool {
	if id == "kolesa" {
		return true
	}

	return false
}
