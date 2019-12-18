package auth_service

type Comparator interface {
	Compare(plaintext string, hashed string) bool
}
