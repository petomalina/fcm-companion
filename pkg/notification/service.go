package notification

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/petomalina/fcm-companion/apis/go-sdk/notification/v1"
	"go.uber.org/zap"
)

type Service struct {
	v1.UnimplementedNotificationServiceServer
	*zap.Logger
}

func (s *Service) PutInstance(ctx context.Context, i *v1.AppInstance) (*empty.Empty, error) {
	panic("implement me")
}

func (s *Service) RemoveToken(ctx context.Context, r *v1.RemoveTokenRequest) (*empty.Empty, error) {
	panic("implement me")
}

func (s *Service) RemoveInstance(ctx context.Context, r *v1.RemoveInstanceRequest) (*empty.Empty, error) {
	panic("implement me")
}

func (s *Service) Send(ctx context.Context, r *v1.SendRequest) (*empty.Empty, error) {
	panic("implement me")
}

func (s *Service) SendAll(ctx context.Context, r *v1.SendAllRequest) (*empty.Empty, error) {
	panic("implement me")
}

func (s *Service) SendMulticast(ctx context.Context, r *v1.SendMulticastRequest) (*empty.Empty, error) {
	panic("implement me")
}

func (s *Service) ListNotifications(ctx context.Context, r *v1.ListNotificationsRequest) (*v1.NotificationList, error) {
	panic("implement me")
}
