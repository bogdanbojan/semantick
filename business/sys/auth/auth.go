// Package auth provides authentication and authorization support.
package auth

import (
	"crypto/rsa"
	"errors"

	"github.com/golang-jwt/jwt"
)

// KeyLookup declares a method set of behavior for looking up
// private and public keys for JWT use.
type KeyLookup interface {
	PrivateKey(kid string) (*rsa.PrivateKey, error)
	PublicKey(kid string) (*rsa.PublicKey, error)
}

// Auth is used to authenticate clients. It can generate a token for a
// set of user claims and recreate the claims by parsing the token.
type Auth struct {
	activeKID string
	keyLookup KeyLookup
	method    jwt.SigningMethod
	keyFunc   func(t *jwt.Token) (interface{}, error)
	parser    jwt.Parser
}

// New creates an Auth to support authentication/authorization.
func New(activeKID string, keyLookup KeyLookup) (*Auth, error) {

	// The activeKID represents the private key used to signed new tokens.
	_, err := keyLookup.PrivateKey(activeKID)
	if err != nil {
		return nil, errors.New("active KID does not exist in store")
	}

	method := jwt.GetSigningMethod("RS256")
	if method == nil {
		return nil, errors.New("configuring algorithm RS256")
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"]
		if !ok {
			return nil, errors.New("missing key id (kid) in token header")
		}
		kidID, ok := kid.(string)
		if !ok {
			return nil, errors.New("user token key id (kid) must be string")
		}
		return keyLookup.PublicKey(kidID)
	}

	// Create the token parser to use. The algorithm used to sign the JWT must be
	// validated to avoid a critical vulnerability:
	// https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
	parser := jwt.Parser{
		ValidMethods: []string{"RS256"},
	}

	a := Auth{
		activeKID: activeKID,
		keyLookup: keyLookup,
		method:    method,
		keyFunc:   keyFunc,
		parser:    parser,
	}

	return &a, nil
}
