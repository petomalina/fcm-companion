package companion

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"go.uber.org/zap"
)

// New returns a new Service with configured firebase services
func New(ctx context.Context, projectID string, logger *zap.Logger) (*Service, error) {
	var firebaseConfig *firebase.Config
	if projectID != "" {
		firebaseConfig = &firebase.Config{ProjectID: projectID}
	}

	firebaseApp, err := firebase.NewApp(context.Background(), firebaseConfig)
	if err != nil {
		logger.Fatal("An error occurred when initializing firebase app", zap.Error(err))
	}

	firestoreClient, err := firebaseApp.Firestore(ctx)
	if err != nil {
		logger.Fatal("An error occurred when initializing firestore", zap.Error(err))
	}

	messagingClient, err := firebaseApp.Messaging(ctx)
	if err != nil {
		logger.Fatal("An error occurred when initializing firebase messaging", zap.Error(err))
	}

	notificationSvc := &Service{
		Logger:          logger,
		FirestoreClient: firestoreClient,
		MessagingClient: messagingClient,
	}

	return notificationSvc, nil
}
