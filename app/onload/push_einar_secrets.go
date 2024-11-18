package onload

import (
	"context"
	"fmt"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	ioc.Registry(pushEinarSecrets)
}

func pushEinarSecrets() {
	secretKeys := []string{
		"GEMINI_API_KEY",
	}

	projectID := "einar-404623"
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	defer client.Close()

	for _, key := range secretKeys {
		secretValue := os.Getenv(key)
		if secretValue == "" {
			log.Printf("Warning: Secret %s is not set", key)
			continue
		}

		secretID := key

		createSecretReq := &secretmanagerpb.CreateSecretRequest{
			Parent:   fmt.Sprintf("projects/%s", projectID),
			SecretId: secretID,
			Secret: &secretmanagerpb.Secret{
				Replication: &secretmanagerpb.Replication{
					Replication: &secretmanagerpb.Replication_Automatic_{
						Automatic: &secretmanagerpb.Replication_Automatic{},
					},
				},
			},
		}

		secret, err := client.CreateSecret(ctx, createSecretReq)
		if err != nil {
			if status.Code(err) == codes.AlreadyExists {
				log.Printf("Secret %s already exists, adding a new version", secretID)
				secret = &secretmanagerpb.Secret{
					Name: fmt.Sprintf("projects/%s/secrets/%s", projectID, secretID),
				}
			} else {
				log.Fatalf("failed to create secret %s: %v", secretID, err)
			}
		}

		addSecretVersionReq := &secretmanagerpb.AddSecretVersionRequest{
			Parent: secret.Name,
			Payload: &secretmanagerpb.SecretPayload{
				Data: []byte(secretValue),
			},
		}

		_, err = client.AddSecretVersion(ctx, addSecretVersionReq)
		if err != nil {
			log.Fatalf("failed to add secret version for %s: %v", secretID, err)
		}

		log.Printf("Secret %s added successfully", secretID)
	}
}
