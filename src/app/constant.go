package app

const (
	ProgName = "GreenViewer"
	Author   = "Carson"
	Version  = "0.0.0"
)

var (
	Port string
)

func init() {
	About = &about{
		ProgName,
		Version,
		Author,
	}
}
