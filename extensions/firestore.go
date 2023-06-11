package extensions

import (
	"sync"

	"cloud.google.com/go/firestore"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
)

var firestore_client *FirestoreClient = &FirestoreClient{}

type FirestoreClient struct {
	cached_client        *firestore.Client
	client_creation_lock sync.Mutex
}

func (fc *FirestoreClient) GetFirestoreClient() *firestore.Client {
	fc.client_creation_lock.Lock()
	defer fc.client_creation_lock.Unlock()

	if fc.cached_client == nil {
		app := firebase_app.getFirestoreApp()
		if app == nil {
			logger.Error("Firebase app is nil")
			return nil
		}

		firestoreClient, err := app.Firestore(firebase_app.ctx)
		if err != nil {
			logger.Error("Failed to create Firestore client", zap.Error(err))
		}
		fc.cached_client = firestoreClient
	}

	return fc.cached_client
}
