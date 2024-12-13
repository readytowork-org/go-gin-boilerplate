package utils

import (
	"os"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

// GetObjectSignedURL the signed url for the stored object
func GetObjectSignedURL(
	bucketName,
	object string,
) (string, error) {

	jsonKey, err := os.ReadFile("serviceAccountKey.json")
	if err != nil {
		return "", nil
	}

	conf, err := google.JWTConfigFromJSON(jsonKey)

	if err != nil {
		return "", err
	}

	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(15 * time.Minute),
	}

	u, err := storage.SignedURL(bucketName, object, opts)

	if err != nil {
		return "", err
	}

	return u, nil
}
