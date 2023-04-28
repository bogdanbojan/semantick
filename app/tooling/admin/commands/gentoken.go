package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/bogdanbojan/semantick/business/sys/auth"
	"github.com/bogdanbojan/semantick/foundation/keystore"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

// GenToken generates a JWT for the specified user.
func GenToken(log *zap.SugaredLogger,  userID string, kid string) error {
	if userID == "" || kid == "" {
		fmt.Println("help: gentoken <user_id> <kid>")
		return nil
	}

	// Construct a key store based on the key files stored in
	// the specified directory.
	keysFolder := "zarf/keys/"
	ks, err := keystore.NewFS(os.DirFS(keysFolder))
	if err != nil {
		return fmt.Errorf("reading keys: %w", err)
	}

	// Init the auth package.
	activeKID := "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"
	a, err := auth.New(activeKID, ks)
	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}

	// Generating a token requires defining a set of claims. In this applications
	// case, we only care about defining the subject and the user in question and
	// the roles they have on the database. This token will expire in a year.
	//
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)
    claims := auth.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "service project",
			Subject:   "foo-bar-hash",
			ExpiresAt: time.Now().Add(8760 * time.Hour).Unix(),
			IssuedAt:  time.Now().UTC().Unix(),
		},
		Roles: []string{"ADMIN"},
	}

	// This will generate a JWT with the claims embedded in them. The database
	// with need to be configured with the information found in the public key
	// file to validate these claims. Dgraph does not support key rotate at
	// this time.
	token, err := a.GenerateToken(claims)
	if err != nil {
		return fmt.Errorf("generating token: %w", err)
	}

	fmt.Printf("-----BEGIN TOKEN-----\n%s\n-----END TOKEN-----\n", token)
	return nil
}

