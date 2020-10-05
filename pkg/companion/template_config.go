package companion

import (
	"cloud.google.com/go/firestore"
	"context"
	v1 "github.com/petomalina/fcm-companion/apis/go-sdk/notification/v1"
	_ "gocloud.dev/runtimevar/gcpruntimeconfig"
	"google.golang.org/api/iterator"
)

const (
	configsCollection = "fcm-companion-configs"
)

func fetchConfig(ctx context.Context, fsClient *firestore.Client, collectionPrefix string) (*v1.NotificationConfig, error) {
	configCol := fsClient.Collection(collectionPrefix + configsCollection)

	config := &v1.NotificationConfig{}

	docs := configCol.Documents(ctx)

	for {
		docSnap, err := docs.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, err
		}

		// get the part of the config into the separate config
		var configPart v1.NotificationConfig
		err = docSnap.DataTo(&configPart)
		if err != nil {
			return nil, err
		}

		// transfer to the main config
		config.Messages = append(config.Messages, configPart.Messages...)
	}

	return config, nil
}
