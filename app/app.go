package app

type about struct {
    ProgName string
    Version  string
    Author   string
}

var About *about

const (
    ProgName = "excel.exe"
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
