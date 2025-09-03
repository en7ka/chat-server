package app

import (
	"context"
	"log"

	"github.com/en7ka/chat-server/internal/api/chat"
	"github.com/en7ka/chat-server/internal/client/db"
	"github.com/en7ka/chat-server/internal/client/db/pg"
	"github.com/en7ka/chat-server/internal/client/db/transaction"
	"github.com/en7ka/chat-server/internal/closer"
	"github.com/en7ka/chat-server/internal/config"
	repoinf "github.com/en7ka/chat-server/internal/repository/chat"
	servinf "github.com/en7ka/chat-server/internal/service/chat"

	userRepo "github.com/en7ka/chat-server/internal/repository/repointerface"
	userService "github.com/en7ka/chat-server/internal/service/servinterface"
)

type serviceProvaider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig

	dbClient       db.Client
	txManager      db.TxManager
	userRepository userRepo.ChatRepository
	userService    userService.ChatService

	userImpl *chat.Controller
}

func newServiceProvider() *serviceProvaider {
	return &serviceProvaider{}
}

func (s *serviceProvaider) GetPGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("NewPGConfig error: %v", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvaider) GetGRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("NewGRPCConfig error: %v", err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvaider) GetHTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("NewHTTPConfig error: %v", err)
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvaider) GetSwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("NewSwaggerConfig error: %v", err)
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig

}
func (s *serviceProvaider) GetDBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.GetPGConfig().DSN())
		if err != nil {
			log.Fatalf("NewDBClient error: %v", err)
		}

		if err = cl.DB().Ping(ctx); err != nil {
			log.Fatalf("NewDBClient error: %v", err)
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvaider) GetTxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.GetDBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvaider) GetUserRepository(ctx context.Context) userRepo.ChatRepository {
	if s.userRepository == nil {
		s.userRepository = repoinf.NewRepository(s.GetDBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvaider) GetUserService(ctx context.Context) userService.ChatService {
	if s.userService == nil {
		s.userService = servinf.NewService(
			s.GetUserRepository(ctx),
			s.GetTxManager(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvaider) GetUserImpl(ctx context.Context) *chat.Controller {
	if s.userImpl == nil {
		s.userImpl = chat.NewImplementation(s.GetUserService(ctx))
	}

	return s.userImpl
}
