package clients

import (
	"context"
	"os"
	"sync"

	socialPb "github.com/Kotlang/notificationGo/generated/social"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// SocialClient encapsulates functionalities related to interacting with the social networking service.
type SocialClient struct {
	cachedConn       *grpc.ClientConn
	connCreationLock sync.Mutex
}

// NewSocialClient creates a new instance of SocialClient.
func NewSocialClient() *SocialClient {
	return &SocialClient{}
}

// GetEventSubscribers retrieves subscribers for a given event asynchronously.
func (c *SocialClient) GetEventSubscribers(grpcContext context.Context, tenant string, eventId string) chan []string {
	subscribers := make(chan []string)

	go func() {
		if eventId == "" {
			logger.Error("Event id is empty")
			return
		}

		conn := c.getConnection()
		if conn == nil {
			logger.Error("Failed to get connection")
			return
		}

		client := socialPb.NewEventsClient(conn)
		ctx := c.prepareCallContext(grpcContext, tenant)
		if ctx == nil {
			logger.Error("Failed to prepare call context")
			return
		}

		response, err := client.GetEventSubscribers(ctx, &socialPb.EventIdRequest{EventId: eventId})
		if err != nil {
			logger.Error("Failed to get event subscribers", zap.Error(err))
			return
		}

		subscribers <- response.GetUserId()
		close(subscribers)
	}()
	return subscribers
}

func (c *SocialClient) getConnection() *grpc.ClientConn {
	c.connCreationLock.Lock()
	defer c.connCreationLock.Unlock()

	if c.cachedConn == nil || c.cachedConn.GetState().String() != "READY" {
		val, ok := os.LookupEnv("SOCIAL_TARGET")
		if !ok || val == "" {
			logger.Error("SOCIAL_TARGET is not set")
			return nil
		}

		conn, err := grpc.Dial(val, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			logger.Error("Failed to connect to social server", zap.Error(err))
			return nil
		}
		c.cachedConn = conn
	}

	return c.cachedConn
}

func (c *SocialClient) prepareCallContext(grpcContext context.Context, tenant string) context.Context {
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
