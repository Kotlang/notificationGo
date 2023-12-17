package extensions

import (
	"context"
	"os"
	"sync"

	generated "github.com/Kotlang/notificationGo/generated"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type SocialClient struct {
	cached_conn        *grpc.ClientConn
	conn_creation_lock sync.Mutex
}

var social_clien *SocialClient = &SocialClient{}

func (c *SocialClient) getConnection() *grpc.ClientConn {
	c.conn_creation_lock.Lock()
	defer c.conn_creation_lock.Unlock()

	if c.cached_conn == nil || c.cached_conn.GetState().String() != "READY" {

		val, ok := os.LookupEnv("SOCIAL_TARGET")
		if !ok || val == "" {
			logger.Error("SOCIAL_TARGET is not set")
		}

		conn, err := grpc.Dial(val, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			logger.Error("Failed to connect to social server", zap.Error(err))
		}
		c.cached_conn = conn
	}

	return c.cached_conn
}

func GetEventSubscribers(grpcContext context.Context, tenant string, eventId string) chan []string {
	subscribers := make(chan []string)

	go func() {
		if eventId == "" {
			logger.Error("Event id is empty")
			return
		}

		conn := social_clien.getConnection()
		if conn == nil {
			logger.Error("Failed to get connection")
			return
		}

		client := generated.NewEventsClient(conn)
		ctx := prepareCallContext(grpcContext, tenant)
		if ctx == nil {
			logger.Error("Failed to prepare call context")
			return
		}

		response, err := client.GetEventSubscribers(ctx, &generated.EventIdRequest{EventId: eventId})
		if err != nil {
			logger.Error("Failed to get event subscribers", zap.Error(err))
			return
		}

		subscribers <- response.GetUserId()
		close(subscribers)
	}()
	return subscribers
}

func prepareCallContext(grpcContext context.Context, tenant string) context.Context {

	// prepare the context
	var md metadata.MD
	if tenant == "neptune" {

		devJWTToken := os.Getenv("DEFAULT_USER_JWT_TOKEN_DEV")
		md = metadata.Pairs("authorization", "bearer "+devJWTToken)
	} else {
		prodJWTToken := os.Getenv("DEFAULT_USER_JWT_TOKEN_PROD")

		md = metadata.Pairs("authorization", "bearer "+prodJWTToken)
	}
	ctx := metadata.NewOutgoingContext(context.TODO(), md)

	return ctx
}
