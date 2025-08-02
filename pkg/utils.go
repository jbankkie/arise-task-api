package pkg

// Utility function to check if a string is empty
func IsEmpty(s string) bool {
    return len(s) == 0
}

// Utility function to trim whitespace from a string
func Trim(s string) string {
    return strings.TrimSpace(s)
}

// Utility function to convert a string to lowercase
func ToLower(s string) string {
    return strings.ToLower(s)
}

// Utility function to check if a string is a valid email format
func IsValidEmail(email string) bool {
    // Simple regex for email validation
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(email)
}