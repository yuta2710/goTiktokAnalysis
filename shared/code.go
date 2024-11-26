package shared

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

type AuthIdProvider interface {
	Encode(userId int) string
	Decode(authKey string) (int, error)
}

// Strategy Pattern for encode/decode auth keys
type Base64AuthIdProvider struct {
}

func (p *Base64AuthIdProvider) Encode(id int, field string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s-%d", field, id)))
}

func (p *Base64AuthIdProvider) Decode(encodedId string, field string) (int, error) {
	fmt.Printf("\nStarting to decoding the ID: %s\n", encodedId)

	// Step 1: Base64 decode
	decoded, err := base64.StdEncoding.DecodeString(encodedId)
	if err != nil {
		return 0, fmt.Errorf("failed to decode id: %v", err)
	}
	fmt.Printf("Decoded string result then: %s\n", string(decoded))

	// Step 2: Split the string into parts
	parts := strings.Split(string(decoded), "-")
	if len(parts) != 2 || parts[0] != field {
		return 0, fmt.Errorf("invalid authId format: %s", string(decoded))
	}
	fmt.Printf("Parts after split: %v\n", parts)

	// Step 3: Convert the ID part to an integer
	decodedId, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("failed to parse id from encoded ID: %v", err)
	}
	fmt.Printf("Parsed id: %d\n", decodedId)

	return decodedId, nil
}
