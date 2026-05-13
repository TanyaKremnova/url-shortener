package service

import (
    "crypto/rand"
    "math/big"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const codeLength = 6

// GenerateShortCode generates a cryptographically random 6-char base62 string
// Example output: "xK3mPq"
func GenerateShortCode() (string, error) {
    code := make([]byte, codeLength)

    for i := range code {
        // Pick a random index into our alphabet
        n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
        if err != nil {
            return "", err
        }
        code[i] = alphabet[n.Int64()]
    }

    return string(code), nil
}