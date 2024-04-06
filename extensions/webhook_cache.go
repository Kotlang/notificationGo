package extensions

import (
	"sync"

	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
)

// Cache struct to manage job cache
type Cache struct {
	data sync.Map
	db   db.NotificationDbInterface
}

type Cacheable interface {
	GetType() string
}

var WhatsappCache *Cache

func GetCache(db db.NotificationDbInterface) *Cache {
	if WhatsappCache == nil {
		WhatsappCache = &Cache{
			db: db,
		}
	}

	return WhatsappCache
}

// Add adds a new item to the cache
func (c *Cache) Add(transID string, item Cacheable) {

	// Load the existing value associated with the key
	existing, _ := c.data.LoadOrStore(transID, []Cacheable{})

	// Append the new item to the existing slice and store it back
	c.data.Store(transID, append(existing.([]Cacheable), item))
}

// UpdateDB updates the database with the cache contents
func (c *Cache) UpdateDB() {
	c.data.Range(func(key, value interface{}) bool {
		transID := key.(string)
		items := value.([]Cacheable)

		// fetch the message from the database
		message := c.db.Message().GetMessageByTransactionId(transID)
		if message == nil {
			logger.Error("Failed to get message by transaction ID", zap.String("transactionID", transID))
			logger.Info("Clearing cache for this key", zap.String("transactionID", transID))
			c.data.Delete(transID)
			return true
		}

		for _, item := range items {

			switch item.GetType() {
			case "delivery":
				it := item.(TrimmedDeliveryInfo)
				switch it.DeliveryStatus {
				case "Delivered":
					message.RecievedBy = append(message.RecievedBy, it.Recipient)
				case "Read":
					message.ReadBy = append(message.ReadBy, it.Recipient)
				case "Failed":
					message.FailedRecipients = append(message.FailedRecipients, it.Recipient)
				default:

				}
			case "buttonReply":
				it := item.(ReplyInfo)

				if message.Responses == nil {
					message.Responses = make(map[string]models.Reply)
				}

				if _, ok := message.Responses[it.Waid]; !ok {
					message.Responses[it.Waid] = models.Reply{
						Replies: []string{},
					}
				}

				message.Responses[it.Waid] = models.Reply{
					Replies: append(message.Responses[it.Waid].Replies, it.ButtonText),
				}
			}
		}

		// update the message in the database
		err := <-c.db.Message().Save(message)
		if err != nil {
			logger.Error("Failed to save message info", zap.Error(err))
			return true // continue without deleting the key from cache
		}

		// Clear the cache for this key after updating the database
		c.data.Delete(transID)

		return true
	})
}
