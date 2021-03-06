Package passwordhash
=====================

	import "github.com/dchest/passwordhash"

Package passwordhash implements safe password hashing and comparison.

Hashes are derived using PBKDF2-HMAC-SHA256 function with 5000 iterations
(by default), 32-byte salt and 64-byte output.

Note: you must not allow users to change parameters of PasswordHash, such as
the number of iterations, directly. If a malicious user can change the
number of iterations, he can set it too high, and it will lead to DoS.

Example usage:

	ph := passwordhash.New("hello, world")
	// Store ph somewhere...
	// Later, when user provides a password:
	if ph.EqualToPassword("hello, world") {
		// Password's okay, user authorized...
	}


Constants
---------

	const (
	    // Default number of iterations for PBKDF2
	    DefaultIterations = 5000
	    // Default salt length
	    SaltLen = 32
	    // Default hash length
	    HashLen = 64
	)



Types
-----

	type PasswordHash struct {
	    Iter int
	    Salt []byte
	    Hash []byte
	}
	
PasswordHash stores hash, salt, and number of iterations.

### func New

	func New(password string) *PasswordHash
	
New returns a new password hash derived from the provided password,
a random salt, and the default number of iterations.
The function causes runtime panic if it fails to get random salt.

### func NewIter

	func NewIter(password string, iter int) *PasswordHash
	
NewIter returns a new password hash derived from the provided password,
the number of iterations, and a random salt.
The function causes runtime panic if it fails to get random salt.

### func NewSaltIter

	func NewSaltIter(password string, salt []byte, iter int) *PasswordHash
	
NewSaltIter creates a new password hash from the provided password, salt,
and the number of iterations.

### func (*PasswordHash) EqualToPassword

	func (ph *PasswordHash) EqualToPassword(password string) bool
	
EqualToPassword returns true if the password hash was derived from the provided password.
This function uses constant time comparison.

### func (*PasswordHash) String

	func (ph *PasswordHash) String() string
	
String returns a string representation of the password hash.

