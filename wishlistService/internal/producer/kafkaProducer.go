package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"pkg/types/trace"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer - продюсер для отправки сообщений в Kafka
type KafkaProducer struct {
	log    *slog.Logger
	writer *kafka.Writer
	topic  string
}

func NewKafkaProducer(log *slog.Logger, broker string, topic string) *KafkaProducer {

	if err := createTopicIfNotExists(log, broker, topic); err != nil {
		log.Warn("failed to create topic:", slog.String("err", err.Error()))
	}

	wrt := &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},   // Равномерное распределение по партициям
		Async:        true,                  // Асинхронный режим
		RequiredAcks: kafka.RequireNone,     // Не ждем подтверждения
		BatchSize:    100,                   // Отправляем пачками до 100 сообщений
		BatchTimeout: 10 * time.Millisecond, // Ждем не больше 10мс
	}

	return &KafkaProducer{
		log:    log,
		writer: wrt,
		topic:  topic,
	}
}

// PublishItemBooked отправляет событие о бронировании
func (p *KafkaProducer) PublishItemBooked(ctx context.Context, event ItemBookedEvent) error {
	const logPrefix = "producer.KafkaProducer.PublishItemBooked"
	log := p.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("PublishItemBooked")

	// Сериализуем событие в JSON
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(event.ItemID), // Ключ для партиционирования
		Value: data,                 // Тело сообщения
		Headers: []kafka.Header{ // Заголовки для фильтрации
			{Key: "event-type", Value: []byte(string(event.EventType))},
			{Key: "event-version", Value: []byte(event.EventVersion)},
			{Key: "item-id", Value: []byte(event.ItemID)},
		},
		Time: time.Now(),
	}

	// асинхронно
	go func() {
		if err := p.writer.WriteMessages(context.Background(), msg); err != nil {
			log.Error("Failed to send booked event to Kafka", slog.String("err", err.Error()))
		} else {
			log.Info("Booked event sent for item", slog.String("item", event.ItemID))
		}
	}()

	return nil
}

// PublishItemUnbooked отправляет событие об отмене бронирования
func (p *KafkaProducer) PublishItemUnbooked(ctx context.Context, event ItemUnbookedEvent) error {
	const logPrefix = "producer.KafkaProducer.PublishItemUnbooked"
	log := p.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("PublishItemUnbooked")

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(event.ItemID),
		Value: data,
		Headers: []kafka.Header{
			{Key: "event-type", Value: []byte(string(event.EventType))},
			{Key: "event-version", Value: []byte(event.EventVersion)},
			{Key: "item-id", Value: []byte(event.ItemID)},
			{Key: "reason", Value: []byte(event.Reason)},
		},
		Time: time.Now(),
	}

	// асинхронно
	go func() {
		if err := p.writer.WriteMessages(context.Background(), msg); err != nil {

			log.Error("Failed to send unbooked event to Kafka", slog.String("err", err.Error()))
		} else {
			log.Info("Unbooked event sent for item", slog.String("item", event.ItemID))
		}
	}()

	return nil
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

func createTopicIfNotExists(l *slog.Logger, broker string, topic string) error {
	const logPrefix = "producer.KafkaProducer.createTopicIfNotExists"
	log := l.With(
		slog.String("where", logPrefix),
	)

	log.Info("createTopicIfNotExists", slog.String("broker", broker), slog.String("topic", topic))

	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return err
	}

	// Ищем наш топик
	for _, p := range partitions {
		if p.Topic == topic {
			log.Info("Topic already exists", slog.String("broker", broker), slog.String("topic", topic))
			return nil
		}
	}

	// Создаем контроллер коннект для создания топиков
	controller, err := conn.Controller()
	if err != nil {
		log.Error("failed to get controller", slog.String("broker", broker), slog.String("topic", topic))
		return fmt.Errorf("failed to get controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Error("failed to dial controller", slog.String("broker", broker), slog.String("topic", topic))
		return fmt.Errorf("failed to dial controller: %w", err)
	}
	defer controllerConn.Close()

	// Конфигурация топика
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		},
	}

	// Создаем топик
	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Error("failed to create topic", slog.String("broker", broker), slog.String("topic", topic))
		return fmt.Errorf("failed to create topic: %w", err)
	}

	log.Info("Topic Created", slog.String("broker", broker), slog.String("topic", topic))
	return nil
}
