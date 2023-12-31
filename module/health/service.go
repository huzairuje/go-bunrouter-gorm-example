package health

import (
	"context"
	"errors"
	"go-bunrouter-gorm-example/module/primitive"

	"github.com/go-redis/redis"
	logger "go-bunrouter-gorm-example/infrastructure/log"
)

type InterfaceService interface {
	CheckUpTime(ctx context.Context) (resp primitive.HealthResp, err error)
}

type Service struct {
	repository  RepositoryInterface
	redisClient *redis.Client
}

func NewService(repository RepositoryInterface, redisClient *redis.Client) InterfaceService {
	return &Service{
		repository:  repository,
		redisClient: redisClient,
	}
}

func (u *Service) CheckUpTime(ctx context.Context) (primitive.HealthResp, error) {
	ctxName := "CheckUpTime"

	if u.repository == nil {
		err := errors.New("repository doesn't initiate on the boot file")
		return primitive.HealthResp{}, err
	}

	if u.redisClient == nil {
		err := errors.New("redis client doesn't initiate on the boot file")
		return primitive.HealthResp{}, err
	}

	errCheckDb := u.repository.CheckUpTimeDB(ctx)
	if errCheckDb != nil {
		logger.Error(ctx, ctxName, "got error when %s : %v", ctxName, errCheckDb)
		return primitive.HealthResp{}, errCheckDb
	}

	errCheckRedis := u.redisClient.Ping().Err()
	if errCheckRedis != nil {
		logger.Error(ctx, ctxName, "got error when %s : %v", ctxName, errCheckRedis)
		return primitive.HealthResp{}, errCheckRedis
	}

	return primitive.HealthResp{
		Db:    "healthy",
		Redis: "healthy",
	}, nil
}
