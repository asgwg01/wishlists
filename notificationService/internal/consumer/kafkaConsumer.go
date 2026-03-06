package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"notificationService/internal/config"
	"notificationService/internal/domain/models"
	email "notificationService/internal/notifier"
	"time"

	"github.com/asgwg01/wishlists/pkg/types/trace"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	log      *slog.Logger
	reader   *kafka.Reader
	notifier *email.EmailNotifier
}

func NewKafkaConsumer(log *slog.Logger, cfg config.KafkaConfig, notifier *email.EmailNotifier) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{cfg.BrokerUrl + ":" + cfg.BrokerPort},
		Topic:          cfg.Topic,
		GroupID:        cfg.GroupId,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		MaxWait:        1 * time.Second,
		CommitInterval: time.Second,
		StartOffset:    kafka.FirstOffset,
	})

	return &KafkaConsumer{
		log:      log,
		reader:   reader,
		notifier: notifier,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) error {
	const logPrefix = "konsumer.KafkaConsumer.Start"
	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("Starting Kafka consumer", slog.String("topic", c.reader.Config().Topic), slog.String("group", c.reader.Config().GroupID))

	c.notifier.Start()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Error("Error reading message", slog.String("err", err.Error()))
				continue
			}

			// в горутине
			go c.processMessage(ctx, msg)
		}
	}
}

func (c *KafkaConsumer) processMessage(ctx context.Context, msg kafka.Message) {
	const logPrefix = "konsumer.KafkaConsumer.processMessage"
	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("Received message",
		slog.Int64("offset", msg.Offset),
		slog.Int("partition", msg.Partition),
		slog.String("key", string(msg.Key)),
	)

	var eventType string
	for _, header := range msg.Headers {
		if header.Key == "event_type" {
			eventType = string(header.Value)
			break
		}
	}

	switch eventType {
	case "item.booked":
		c.processBookedEvent(ctx, msg.Value)
	case "item.unbooked":
		c.processUnbookedEvent(ctx, msg.Value)
	default:
		log.Warn("Unknown event type",
			slog.String("type", eventType),
		)
	}
}

func (c *KafkaConsumer) processBookedEvent(ctx context.Context, data []byte) {
	const logPrefix = "konsumer.KafkaConsumer.processBookedEvent"
	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	var event models.ItemBookedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		log.Error("Failed to parse booked event", slog.String("err", err.Error()))
		return
	}

	log.Info("Processing booked event",
		slog.String("event_id", event.ItemID),
		slog.String("item_name", event.ItemName),
	)

	// Отправляем email владельцу
	subject := fmt.Sprintf("Твоя хотелка '%s' была забронирована!", event.ItemName)
	body := fmt.Sprintf(`Здарова %s!
	Кое кто хочет сделать тебе подарок!
	%s (%s) Забронировал позицию "%s" из твоего вишлиста "%s".
	Есть вероятность, что, вскоре тебе достанется подарок)`,
		event.OwnerName,
		event.BookedByName,
		event.BookedByEmail,
		event.ItemName,
		event.WishlistName)

	if err := c.notifier.SendEmail(event.OwnerEmail, subject, body); err != nil {
		log.Error("Failed to send email to owner", slog.String("err", err.Error()))
	}

	// Отправляем подтверждение тому, кто забронировал
	subject = fmt.Sprintf("Вы забронировали '%s'", event.ItemName)
	body = fmt.Sprintf(`Приветик %s!
	Вы успешно забронировали позицию "%s" пользователя %s',
	в его вишлисте "%s".
	Владелец будет извещен о вашем бронировании!
	Уверен подарок будет супер!`,
		event.BookedByName,
		event.ItemName,
		event.OwnerName,
		event.WishlistName)

	if err := c.notifier.SendEmail(event.BookedByEmail, subject, body); err != nil {
		log.Error("Failed to send email to booker", slog.String("err", err.Error()))
	}
}

func (c *KafkaConsumer) processUnbookedEvent(ctx context.Context, data []byte) {
	const logPrefix = "konsumer.KafkaConsumer.processUnbookedEvent"
	log := c.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	var event models.ItemUnbookedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		log.Error("Failed to parse unbooked event", slog.String("err", err.Error()))
		return
	}

	log.Info("Processing booked event",
		slog.String("event_id", event.ItemID),
		slog.String("item_name", event.ItemName),
		slog.String("reason", event.Reason),
	)

	switch event.Reason {
	case "cancelled_by_user":
		// Пользователь передумал дарить
		subject := fmt.Sprintf("Бронирование отменено!: '%s'", event.ItemName)
		body := fmt.Sprintf(`Привет %s,
		Пользователь %s (%s) снял бронь с позиции "%s" из вишлиста "%s".
		Может еще передумает)`,
			event.OwnerName,
			event.UnbookedByName,
			event.UnbookedByEmail,
			event.ItemName,
			event.WishlistName)

		if err := c.notifier.SendEmail(event.OwnerEmail, subject, body); err != nil {
			log.Error("Failed to send email", slog.String("err", err.Error()))
		}

	case "cancelled_by_owner":
		// Владелец отменил бронирование
		subject := fmt.Sprintf("Бронирование было отменено владельцем! '%s'", event.ItemName)
		body := fmt.Sprintf(`Приветик %s!
		Владелец вишлиста "%s",  снял бронь у забронированной, тобой позиции "%s".
		Пиши - узнавай`,
			event.BookedByName,
			event.WishlistName,
			event.ItemName)

		if err := c.notifier.SendEmail(event.BookedByEmail, subject, body); err != nil {
			log.Error("Failed to send email", slog.String("err", err.Error()))
		}
	}
}

func (c *KafkaConsumer) Close() error {
	c.notifier.Stop()
	return c.reader.Close()
}
