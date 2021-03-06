// Package passwordhash implements safe password hashing and comparison.
//
// Hashes are derived using PBKDF2-HMAC-SHA256 function with 5000 iterations
// (by default), 32-byte salt and 64-byte output.
//
// Note: you must not allow users to change parameters of PasswordHash, such as
// the number of iterations, directly. If a malicious user can change the
// number of iterations, he can set it too high, and it will lead to DoS.
//
// Example usage:
//
//	ph := passwordhash.New("hello, world")
//	// Store ph somewhere...
//	// Later, when user provides a password:
//	if ph.EqualToPassword("hello, world") {
//		// Password's okay, user authorized...
//	}
//
package passwordhash

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"libs/pbkdf2"
	"io"
)

// PasswordHash stores hash, salt, and number of iterations.
type PasswordHash struct {
	Iter int
	Salt []byte
	Hash []byte
}

const (
	// Default number of iterations for PBKDF2
	DefaultIterations = 5000
	// Default salt length
	SaltLen = 32
	// Default hash length
	HashLen = 64
)

// getSalt returns a new random salt.
// The function causes runtime panic if it fails to read from random source.
func getSalt() []byte {
	salt := make([]byte, SaltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic("error reading from random source: " + err.String())
	}
	return salt
}

// New returns a new password hash derived from the provided password, 
// a random salt, and the default number of iterations.
// The function causes runtime panic if it fails to get random salt.
func New(password string) *PasswordHash {
	return NewSaltIter(password, getSalt(), DefaultIterations)
}

// NewIter returns a new password hash derived from the provided password,
// the number of iterations, and a random salt.
// The function causes runtime panic if it fails to get random salt.
func NewIter(password string, iter int) *PasswordHash {
	return NewSaltIter(password, getSalt(), iter)
}

// NewSaltIter creates a new password hash from the provided password, salt,
// and the number of iterations.
func NewSaltIter(password string, salt []byte, iter int) *PasswordHash {
	return &PasswordHash{iter, salt,
		pbkdf2.WithHMAC(sha256.New, []byte(password), salt, iter, HashLen)}
}

// EqualToPassword returns true if the password hash was derived from the provided password.
// This function uses constant time comparison.
func (ph *PasswordHash) EqualToPassword(password string) bool {
	provided := NewSaltIter(password, ph.Salt, ph.Iter)
	if len(ph.Hash) != len(provided.Hash) {
		return false
	}
	return subtle.ConstantTimeCompare(ph.Hash, provided.Hash) == 1
}

// String returns a string representation of the password hash.
func (ph *PasswordHash) String() string {
	return fmt.Sprintf("&PasswordHash{Iter: %d, Salt: %x, Hash: %x}",
		ph.Iter, ph.Salt, ph.Hash)
}
