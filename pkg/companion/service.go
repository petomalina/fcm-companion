package companion

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/v4/messaging"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/petomalina/fcm-companion/apis/go-sdk/notification/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	instancesCollection = "fcm-companion-instances"
)

// Service is the implementation of the Notification API
type Service struct {
	v1.UnimplementedNotificationServiceServer
	*zap.Logger

	FirestoreClient *firestore.Client
	MessagingClient *messaging.Client

	// CollectionPrefix is used to define Firestore collection prefixes.
	// This is useful when there would be conflict with already existing collections
	CollectionPrefix string
}

// Register registers this service to the provided grpc server
func (s *Service) Register(server *grpc.Server) {
	v1.RegisterNotificationServiceServer(server, s)
}

// RegisterGateway registers this service to the provided http mux
func (s *Service) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, bind string, opts []grpc.DialOption) error {
	return v1.RegisterNotificationServiceHandlerFromEndpoint(ctx, mux, bind, opts)
}

func (s *Service) PutInstance(ctx context.Context, i *v1.AppInstance) (*empty.Empty, error) {
	if err := i.Validate(); err != nil {
		return &empty.Empty{}, err
	}

	// get the document reference (this won't read it)
	doc := s.FirestoreClient.Doc(s.CollectionPrefix + instancesCollection + "/" + i.InstanceId)

	// we don't need to read the document if we can just merge it
	if len(i.Labels) <= 0 {
		_, err := doc.Set(ctx, i, firestore.MergeAll)
		return &empty.Empty{}, err
	}

	// overwrite the labels if set, with the whole doc
	_, err := doc.Set(ctx, i)
	return &empty.Empty{}, err
}

func (s *Service) RemoveToken(ctx context.Context, r *v1.RemoveTokenRequest) (*empty.Empty, error) {
	if err := r.Validate(); err != nil {
		return &empty.Empty{}, err
	}

	doc := s.FirestoreClient.Doc(s.CollectionPrefix + instancesCollection + "/" + r.InstanceId)

	// this only resets the value of the token for the instance ID
	_, err := doc.Update(ctx, []firestore.Update{
		{
			Path: "token", Value: "",
		},
	})

	return &empty.Empty{}, err
}

func (s *Service) RemoveInstance(ctx context.Context, r *v1.RemoveInstanceRequest) (*empty.Empty, error) {
	if err := r.Validate(); err != nil {
		return &empty.Empty{}, err
	}

	doc := s.FirestoreClient.Doc(s.CollectionPrefix + instancesCollection + "/" + r.InstanceId)

	_, err := doc.Delete(ctx)
	return &empty.Empty{}, err
}

func (s *Service) Send(ctx context.Context, r *v1.SendRequest) (*empty.Empty, error) {
	if err := r.Validate(); err != nil {
		return &empty.Empty{}, err
	}

	panic("implement me")
}

func (s *Service) SendAll(ctx context.Context, r *v1.SendAllRequest) (*empty.Empty, error) {
	if err := r.Validate(); err != nil {
		return &empty.Empty{}, err
	}

	panic("implement me")
}

func (s *Service) SendMulticast(ctx context.Context, r *v1.SendMulticastRequest) (*empty.Empty, error) {
	if err := r.Validate(); err != nil {
		return &empty.Empty{}, err
	}

	panic("implement me")
}

func (s *Service) ListNotifications(ctx context.Context, r *v1.ListNotificationsRequest) (*v1.NotificationList, error) {
	if err := r.Validate(); err != nil {
		return &v1.NotificationList{}, err
	}

	panic("implement me")
}
