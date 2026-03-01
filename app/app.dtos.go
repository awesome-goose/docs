package app

// BuildDto contains options for the build command
type BuildDto struct {
	Input   string `flag:"input" default:"docs"`
	Output  string `flag:"output" default:"dist"`
	Title   string `flag:"title" default:"Goose Documentation"`
	BaseURL string `flag:"base-url" default:"/"`
}

// ServeDto contains options for the serve command
type ServeDto struct {
	Dir  string `flag:"dir" default:"dist"`
	Port string `flag:"port" default:"3000"`
}

// HelpDto is an empty DTO for the help command
type HelpDto struct{}
