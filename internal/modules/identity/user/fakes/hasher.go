package fakes

// Hasher is a reversible stand-in for bcrypt. Tests only need hash != plain.
type Hasher struct{}

func (Hasher) Hash(plain string) (string, error) {
	return "hashed:" + plain, nil
}

func (Hasher) Verify(hash, plain string) bool {
	return hash == "hashed:"+plain
}
