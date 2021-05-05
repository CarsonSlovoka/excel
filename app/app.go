package app

type about struct {
    ProgName string
    Version  string
    Author   string
}

var About *about

const (
    ProgName = "GreenViewer"
    Author   = "Carson"
    Version  = "0.0.0"
)

func init() {
    About = &about{
        ProgName,
        Author,
        Version,
    }
}
