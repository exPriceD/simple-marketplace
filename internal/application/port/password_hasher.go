package port

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Verify(hash, plain string) bool
}
