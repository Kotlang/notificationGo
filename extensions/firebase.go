package extensions

import (
	"context"
	"os"
	"sync"

	firebase "firebase.google.com/go/v4"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

var firebase_app *FirebaseApp = &FirebaseApp{}

type FirebaseApp struct {
	cached_app        *firebase.App
	app_creation_lock sync.Mutex
	ctx               context.Context
}

func (fa *FirebaseApp) getFirestoreApp() *firebase.App {
	fa.app_creation_lock.Lock()
	defer fa.app_creation_lock.Unlock()

	if fa.cached_app == nil {
		opt := option.WithCredentialsJSON([]byte(os.Getenv("FCM-TOKEN")))

		ctx := context.Background()

		app, err := firebase.NewApp(ctx, nil, opt)

		if err != nil {
			logger.Error("Failed to initialize Firebase app", zap.Error(err))
			return nil
		}
		fa.ctx = ctx
		fa.cached_app = app
	}

	return fa.cached_app
}
