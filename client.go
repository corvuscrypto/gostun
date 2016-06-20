package gostun

//Client is the main struct that hold the Client info for P2P tx
type Client struct {
	UID        string
	Session    string
	IPAddr     string
	ReturnPort int
	Info       map[string]interface{}
}

//MakeClient takes data packet (eventually) and returns a Client object
func MakeClient() *Client {
	return new(Client)
}
