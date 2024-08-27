package replacer

import (
	"regexp"
	"strings"

	bitwarden "github.com/bitwarden/sdk-go"
)

const (
	lookupRegex = `<bw:(?P<secret>[^ ]+)>`
)

type Replacer struct {
	client LookupInterface
}

type LookupInterface interface {
	Get(secretID string) (*bitwarden.SecretResponse, error)
}

func New(apiURL, identityURL, token string) (*Replacer, error) {
	client, err := bitwarden.NewBitwardenClient(
		&apiURL,
		&identityURL,
	)
	if err != nil {
		return nil, err
	}

	err = client.AccessTokenLogin(token, nil)
	if err != nil {
		return nil, err
	}

	return &Replacer{
		client: client.Secrets(),
	}, nil
}

func (r Replacer) Replace(raw string) (string, error) {

	re, err := regexp.Compile(lookupRegex)
	if err != nil {
		return "", err
	}

	if !re.MatchString(raw) {
		return raw, nil
	}

	result := strings.Clone(raw)
	entries := re.FindAllStringSubmatch(raw, -1)

	for _, entry := range entries {
		path := entry[0]
		secret := entry[1]
		replacement, err := r.client.Get(secret)
		if err != nil {
			return "", err
		}

		result = strings.ReplaceAll(result, path, replacement.Value)
	}

	return result, nil
}
