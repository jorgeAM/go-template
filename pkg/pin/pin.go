package pin

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func Generate4DigitCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%04d", n.Int64()), nil
}
