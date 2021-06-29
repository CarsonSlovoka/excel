package app

const (
    ProgName = "GreenViewer"
    Author   = "Carson"
    Version  = "0.0.0"
    Port     = "7121"
)

func init() {
    About = &about{
        ProgName,
        Version,
        Author,
    }
}
