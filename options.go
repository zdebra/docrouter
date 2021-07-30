package docrouter

type Options struct {
	Title   string
	Version string

	// ServerURLs are used purely for generating OpenAPI schema.
	// It doesn't have any effect on a request host matching.
	Servers []ServerDoc
}

type ServerDoc struct {
	URL         string
	Description string
}

var DefaultOptions = Options{
	Title:   "Default Title",
	Version: "1.0",
	Servers: []ServerDoc{
		{"https://www.example.com/v3", "Production environment API"},
		{"https://test.example.com/v3", "Test environment API"},
	},
}
