package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/storage"
)

func main() {
	store := storage.NewExampleStore()

	// Configure the strategy and key
	config := &fosite.Config{
		AccessTokenLifespan: time.Hour,
		IDTokenLifespan:     time.Hour,
		GlobalSecret:        []byte("some-cool-secret-that-is-32bytes"),
		// RefreshTokenScopes:  []string{}, // or set accessRequest.GrantScope("offline")
	}

	// privateKey is used to sign JWT tokens. The default strategy uses RS256 (RSA Signature with SHA-256)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Create a new Fosite instance
	oauth2Provider := compose.ComposeAllEnabled(config, store, privateKey)

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse request
		mySessionData := new(fosite.DefaultSession)
		accessRequest, err := oauth2Provider.NewAccessRequest(ctx, r, mySessionData)
		if err != nil {
			oauth2Provider.WriteAccessError(ctx, w, accessRequest, err)
			return
		}

		// Create access token
		response, err := oauth2Provider.NewAccessResponse(ctx, accessRequest)
		if err != nil {
			oauth2Provider.WriteAccessError(ctx, w, accessRequest, err)
			return
		}
		log.Printf("%+#v", response)

		// Send response
		oauth2Provider.WriteAccessResponse(ctx, w, accessRequest, response)
	})

	// Start HTTP server
	log.Println("Server is running at http://localhost:3847")
	fmt.Println()
	log.Println(`You can test with

token=$(curl -s -X POST http://localhost:3847/token \
	-H "Content-Type: application/x-www-form-urlencoded" \
	-d "grant_type=password" \
	-d "client_id=my-client" \
	-d "client_secret=foobar" \
	-d "username=peter" \
	-d "password=secret" \
	-d "scope=fosite offline")
refresh_token=$(echo $token | jq -r '.refresh_token')

echo ""
echo "Token:" $token
echo ""
echo "RefreshToken:" $refresh_token

token=$(curl -s -X POST http://localhost:3847/token \
	-H "Content-Type: application/x-www-form-urlencoded" \
	-d "grant_type=refresh_token" \
	-d "client_id=my-client" \
	-d "client_secret=foobar" \
	-d "refresh_token=$refresh_token")
echo ""
echo "NewToken:" $token
	`)
	log.Fatal(http.ListenAndServe(":3847", nil))
}
