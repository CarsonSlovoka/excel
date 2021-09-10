package github

// f'https://api.github.com/repos/{owner}/{repo}/contents/{path}?ref={branch}'
type QueryInfo struct {
    Owner  string
    Repo   string
    Path   string
    Branch string
    Token  string
}
