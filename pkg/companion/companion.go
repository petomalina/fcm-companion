package companion

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"go.uber.org/zap"
)

// New returns a new Service with configured firebase services
func New(ctx context.Context, projectID string, logger *zap.Logger, collectionPrefix string) (*Service, error) {
	var firebaseConfig *firebase.Config
	if projectID != "" {
		firebaseConfig = &firebase.Config{
			ProjectID: projectID,
		}
	}

	firebaseApp, err := firebase.NewApp(context.Background(), firebaseConfig)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	firestoreClient, err := firebaseApp.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing firestore: %v", err)
	}

	messagingClient, err := firebaseApp.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("error firebase messaging: %v", err)
	}

	config, err := fetchConfig(ctx, firestoreClient, collectionPrefix)
	if err != nil {
		return nil, fmt.Errorf("error fetching configuration: %v", err)
	}

	fmt.Println(config.Messages)

	notificationSvc := &Service{
		Logger:           logger,
		FirestoreClient:  firestoreClient,
		MessagingClient:  messagingClient,
		CollectionPrefix: collectionPrefix,
		config:           config,
	}

	return notificationSvc, nil
}
