package docserver

type Options struct {
	Title      string
	Version    string
	ServerURLs []string
}

var DefaultOptions = Options{
	Title:   "Default Title",
	Version: "1.0",
	ServerURLs: []string{
		"http://localhost:8000",
	},
}
