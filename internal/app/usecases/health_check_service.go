package usecases

import "github.com/mirfnsyh/base-service/internal/app/commons"

type IHealthCheckService interface {
	Ping() string
}

type HealthCheckService struct {
	opt commons.Options
}

func NewHealthCheck(opt commons.Options) IHealthCheckService {
	return &HealthCheckService{
		opt: opt,
	}
}

func (svc *HealthCheckService) Ping() string {
	return "PONG"
}
