package utils

func In(label string, set []string) bool {
    for _, val := range set {
        if label == val {
            return true
        }
    }
    return false
}
