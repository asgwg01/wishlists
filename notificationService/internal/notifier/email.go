package email

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
	"notificationService/internal/config"
	"time"
)

type EmailMessage struct {
	To      string
	Subject string
	Body    string
}

type EmailNotifier struct {
	log      *slog.Logger
	host     string
	port     string
	user     string
	password string
	from     string
	isDev    bool

	// пришлось вкрячить, потому что гугл и яндекс на**й послают
	// с smtp.yandex.ru и smtp.gmail.com на 587м порту. Спамер получается)
	// Пришлось использовать
	// mailtrap.io который на всех портах кроме 2525 тоже посылает,
	// а еще для бесплатного использования требует не чаще 1 сообщения в 10 сек
	queue         chan EmailMessage
	rateLimit     time.Duration // интервал между отправками, ставим чуть более 10 секунд
	ctx           context.Context
	cancelCalback context.CancelFunc
}

func NewEmailNotifier(log *slog.Logger, cfg *config.Config) *EmailNotifier {
	return &EmailNotifier{
		log:      log,
		host:     cfg.SMTPConfig.Host,
		port:     cfg.SMTPConfig.Port,
		user:     cfg.SMTPConfig.User,
		password: cfg.SMTPConfig.Password,
		from:     cfg.SMTPConfig.From,
		isDev:    cfg.Env == "local",

		queue:     make(chan EmailMessage, 100),
		rateLimit: 10100 * time.Millisecond,
	}
}

func (n *EmailNotifier) Start() {
	const logPrefix = "notifier.EmailNotifier.Start"
	log := n.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Start")

	ctx, cancel := context.WithCancel(context.Background())
	n.ctx = ctx
	n.cancelCalback = cancel
	go n.loop()
}

func (n *EmailNotifier) Stop() {
	const logPrefix = "notifier.EmailNotifier.Stop"
	log := n.log.With(
		slog.String("where", logPrefix),
	)

	log.Info("Stop")

	n.cancelCalback()
}

func (n *EmailNotifier) loop() {
	const logPrefix = "notifier.EmailNotifier.loop"
	log := n.log.With(
		slog.String("where", logPrefix),
	)

	ticker := time.NewTicker(n.rateLimit)
	defer ticker.Stop()

	for {
		select {
		case <-n.ctx.Done():
			log.Info("Stop EmailNotifier loop")
			return
		case <-ticker.C:
			select {
			case msg := <-n.queue:
				n.sendEmail(msg.To, msg.Subject, msg.Body)
			default:
				// next loop
			}
		}
	}
}

func (n *EmailNotifier) AddToQueue(to, subject, body string) {
	const logPrefix = "notifier.EmailNotifier.AddToQueue"
	log := n.log.With(
		slog.String("where", logPrefix),
	)

	select {
	case n.queue <- EmailMessage{To: to, Subject: subject, Body: body}:
		log.Info("Email queued", slog.String("to", to))
	default:
		log.Info("Queue is full, drop email", slog.String("to", to))
	}
}

func (n *EmailNotifier) SendEmail(to, subject, body string) error {
	n.AddToQueue(to, subject, body)
	return nil
}

func (n *EmailNotifier) sendEmail(to, subject, body string) {
	const logPrefix = "notifier.EmailNotifier.sendEmail"
	log := n.log.With(
		slog.String("where", logPrefix),
	)
	// В режиме разработки просто логируем
	if n.isDev || n.user == "" || n.password == "" {
		log.Info("STUB! Send email",
			slog.String("to", to),
			slog.String("subject", subject),
			slog.String("body", body),
		)
		return
	}

	// RFC 822
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n"+
		"\r\n"+
		"%s\r\n", to, n.from, subject, body))


	// smtp auth
	auth := smtp.PlainAuth("", n.user, n.password, n.host)
	addr := fmt.Sprintf("%s:%s", n.host, n.port)


	err := smtp.SendMail(addr, auth, n.from, []string{to}, msg)
	if err != nil {
		log.Error("failed to send email", slog.String("err", err.Error()))
		return 
	}

	log.Info("Email sent", slog.String("to", to))
}
