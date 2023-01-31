package anh

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type TokenProvider interface {
	GenerateSasToken(audienceUri string, expiry time.Time) (string, time.Time, error)
}

type sasTokenProvider struct {
	keyName string
	key     []byte
}

func NewSasTokenProvider(keyName string, key string) TokenProvider {
	return &sasTokenProvider{
		keyName: keyName,
		key:     []byte(key),
	}
}

func (p *sasTokenProvider) GenerateSasToken(audienceUri string, expiry time.Time) (string, time.Time, error) {
	audienceUri = url.QueryEscape(audienceUri)
	audienceUri = strings.ToLower(audienceUri)

	bytesToSign := []byte(fmt.Sprintf("%s\n%d", audienceUri, expiry.Unix()))

	hmacHash := hmac.New(sha256.New, []byte(p.key))
	hmacHash.Write(bytesToSign)
	signature := hmacHash.Sum(nil)

	signatureString := base64.StdEncoding.EncodeToString(signature)

	token := fmt.Sprintf("SharedAccessSignature sr=%s&sig=%s&se=%d&skn=%s", audienceUri, url.QueryEscape(signatureString), expiry.Unix(), p.keyName)
	return token, expiry, nil
}
