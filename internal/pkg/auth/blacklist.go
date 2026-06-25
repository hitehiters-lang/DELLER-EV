package pkgauth

import "time"

var TokenBlacklist = make(map[string]time.Time)

func CleanupBlacklist() {
	for token, expiresAt := range TokenBlacklist {
		if time.Now().After(expiresAt) {
			delete(TokenBlacklist, token)
		}
	}
}
