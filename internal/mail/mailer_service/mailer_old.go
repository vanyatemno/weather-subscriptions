package mailer_service

//
//import (
//	"context"
//	"fmt"
//	"net/smtp"
//	"strings"
//	"sync"
//	"time"
//	"weather-subscriptions/internal/config"
//)
//
//// Mailer defines the interface for a mail sending service.
//// Implementations are responsible for the mechanics of dispatching emails.
//type Mailer interface {
//	// Send attempts to queue an email message for delivery.
//	// It should be safe for concurrent use.
//	// Returns an error if the message cannot be queued (e.g., service is shutting down,
//	// queue is full, or message validation fails).
//	Send(message MailMessage) error
//
//	// Close gracefully shuts down the mailer_service.
//	// It should ensure that pending tasks are completed or handled according to the
//	// implementation's policy (e.g., attempt to send, log, or discard).
//	// After Close returns, the Mailer should not accept new tasks.
//	// Close should be idempotent.
//	Close() error
//}
//
//// DefaultMailer is a robust implementation of the Mailer interface using SMTP.
//// It employs a pool of worker goroutines to send emails asynchronously, supports retries,
//// and integrates with Go's context mechanism for graceful shutdowns.
//type DefaultMailer struct {
//	smtpConfig     SMTPConfig
//	taskQueue      chan MailMessage   // Channel for queueing email messages.
//	workerPoolSize int                // Number of concurrent worker goroutines.
//	maxRetries     int                // Maximum number of retry attempts for sending an email.
//	retryDelay     time.Duration      // Delay between retry attempts.
//	wg             sync.WaitGroup     // WaitGroup to synchronize worker goroutine shutdown.
//	ctx            context.Context    // Internal context managing the lifecycle of the mailer_service.
//	cancel         context.CancelFunc // Function to cancel the internal context, triggering shutdown.
//	// logger         *log.Logger      // Optional: For logging operational messages.
//}
//
//// NewDefaultMailer creates and initializes a new DefaultMailer instance.
////
//// Parameters:
////   - appCtx: The parent application context. The mailer_service will listen for cancellation of this
////     context to initiate a graceful shutdown. Must not be nil.
////   - cfg: The config of the running application.
////
//// Returns:
////
////	A pointer to the initialized DefaultMailer and an error if initialization fails
////	(e.g., due to invalid configuration parameters or a nil appCtx).
////
//// The mailer_service starts its worker pool upon creation and is ready to send emails immediately.
//// It's the caller's responsibility to call Close() on the mailer_service when it's no longer needed
//// to ensure all resources are released and pending emails are processed according to policy.
//// func NewDefaultMailer(appCtx context.Context, smtpCfg SMTPConfig, workerPoolSize int, queueCapacity int, maxRetries int, retryDelay time.Duration) (*DefaultMailer, error) {
//func NewDefaultMailer(appCtx context.Context, cfg *config.Config) (*DefaultMailer, error) {
//	if appCtx == nil {
//		return nil, fmt.Errorf("application context (appCtx) cannot be nil")
//	}
//	if cfg.Mailer.WorkerPoolSize <= 0 {
//		return nil, fmt.Errorf("workerPoolSize must be greater than 0")
//	}
//	if cfg.Mailer.QueueCapacity <= 0 {
//		return nil, fmt.Errorf("queueCapacity must be greater than 0")
//	}
//	if cfg.Mailer.Host == "" || cfg.Mailer.Port == 0 || cfg.Mailer.From == "" {
//		return nil, fmt.Errorf("SMTP host, port, and from address must be configured")
//	}
//
//	// Create an internal, cancellable context derived from the application context.
//	// This allows the mailer_service to be shut down in two ways:
//	// 1. If appCtx is cancelled (e.g., application is shutting down).
//	// 2. If the mailer_service's Close() method is explicitly called.
//	internalCtx, internalCancel := context.WithCancel(appCtx)
//
//	mailer := &DefaultMailer{
//		smtpConfig: SMTPConfig{
//			Host:     cfg.Mailer.Host,
//			Port:     cfg.Mailer.Port,
//			Username: cfg.Mailer.Username,
//			Password: cfg.Mailer.Password,
//			From:     cfg.Mailer.From,
//		},
//		taskQueue:      make(chan MailMessage, cfg.Mailer.QueueCapacity),
//		workerPoolSize: cfg.Mailer.WorkerPoolSize,
//		maxRetries:     cfg.Mailer.MaxRetries,
//		retryDelay:     time.Duration(float64(cfg.Mailer.RetryDelay) * float64(time.Second)),
//		ctx:            internalCtx,
//		cancel:         internalCancel,
//		// logger:      log.New(os.Stdout, "[DefaultMailer] ", log.LstdFlags), // Example logger setup
//	}
//
//	// Start worker goroutines.
//	mailer.wg.Add(cfg.Mailer.WorkerPoolSize)
//	for i := 0; i < cfg.Mailer.WorkerPoolSize; i++ {
//		go mailer.worker(i)
//	}
//
//	return mailer, nil
//}
//
//// Send queues an email message for asynchronous delivery by a worker.
////
//// Parameters:
////   - message: The MailMessage to be sent. Basic validation (To, Subject, Body) is performed.
////
//// Returns:
////
////	An error if the message validation fails, if the mailer_service is shutting down (context cancelled),
////	or if the task queue is full and the message cannot be queued within a short timeout (1 second).
////
//// This method is safe for concurrent use.
//func (m *DefaultMailer) Send(message MailMessage) error {
//	if len(message.To) == 0 {
//		return fmt.Errorf("message must have at least one recipient in 'To' field")
//	}
//	if message.Subject == "" {
//		return fmt.Errorf("message subject cannot be empty")
//	}
//	if message.Body == "" {
//		return fmt.Errorf("message body cannot be empty")
//	}
//
//	// Attempt to queue the message, respecting context cancellation and a timeout.
//	select {
//	case m.taskQueue <- message:
//		// if m.logger != nil { m.logger.Printf("Email to %v queued successfully.", message.To) }
//		return nil
//	case <-time.After(1 * time.Second): // Timeout for queuing if the queue is full.
//		// if m.logger != nil { m.logger.Printf("Failed to queue email to %v: task queue full or unresponsive.", message.To) }
//		return fmt.Errorf("failed to queue email: task queue full or unresponsive")
//	case <-m.ctx.Done(): // Check if the mailer_service's context is cancelled.
//		// if m.logger != nil { m.logger.Printf("Failed to queue email to %v: mailer_service is shutting down (%v).", message.To, m.ctx.Err()) }
//		return fmt.Errorf("mailer_service is shutting down or context cancelled, cannot queue new email: %w", m.ctx.Err())
//	}
//}
//
//// worker is a long-running goroutine that dequeues email messages from taskQueue
//// and processes them for sending via sendEmailWithRetries.
//// Each worker operates independently. Workers exit when the mailer_service's context (m.ctx)
//// is cancelled or when the taskQueue is closed and subsequently emptied.
//func (m *DefaultMailer) worker(id int) {
//	defer m.wg.Done()
//	// if m.logger != nil { m.logger.Printf("Worker %d started.", id) }
//
//	for {
//		select {
//		case message, ok := <-m.taskQueue:
//			if !ok {
//				// taskQueue has been closed, and it's empty. Worker should exit.
//				// if m.logger != nil { m.logger.Printf("Worker %d: taskQueue closed, exiting.", id) }
//				return
//			}
//			// if m.logger != nil { m.logger.Printf("Worker %d: processing email to %v.", id, message.To) }
//			m.sendEmailWithRetries(message) // This method also respects m.ctx.
//		case <-m.ctx.Done():
//			// Mailer's context has been cancelled. Worker should exit.
//			// if m.logger != nil { m.logger.Printf("Worker %d: context cancelled (%v), exiting.", id, m.ctx.Err()) }
//			return
//		}
//	}
//}
//
//// sendEmailWithRetries attempts to send an email message, with retries on failure, up to m.maxRetries.
//// It respects the mailer_service's context (m.ctx) and will stop retrying if the context is cancelled.
//func (m *DefaultMailer) sendEmailWithRetries(message MailMessage) {
//	var err error
//	for i := 0; i < m.maxRetries; i++ {
//		// Check context before attempting to send, especially important for retries.
//		select {
//		case <-m.ctx.Done():
//			// if m.logger != nil { m.logger.Printf("Context cancelled before sending/retrying email to %v (%v). Email not sent.", message.To, m.ctx.Err()) }
//			return // Context cancelled, do not attempt to send or retry.
//		default:
//			// Context is not (yet) cancelled, proceed with send attempt.
//		}
//
//		err = m.sendDirect(message)
//		if err == nil {
//			// if m.logger != nil { m.logger.Printf("Email to %v sent successfully on attempt %d.", message.To, i+1) }
//			return // Successfully sent.
//		}
//		// if m.logger != nil { m.logger.Printf("Attempt %d to send email to %v failed: %v. Retrying in %v...", i+1, message.To, err, m.retryDelay) }
//
//		// Wait for retryDelay or until context is cancelled.
//		select {
//		case <-time.After(m.retryDelay):
//			// Continue to the next retry iteration.
//		case <-m.ctx.Done():
//			// if m.logger != nil { m.logger.Printf("Context cancelled during retry delay for email to %v (%v). Email not sent.", message.To, m.ctx.Err()) }
//			// Log or handle that the message was not sent due to shutdown/cancellation.
//			return
//		}
//	}
//	// If all retries fail.
//	// if m.logger != nil { m.logger.Printf("Failed to send email to %v after %d retries. Last error: %v", message.To, m.maxRetries, err) }
//}
//
//// sendDirect performs the actual SMTP sending of a single email message.
//// This method is responsible for constructing the email (headers, body) and
//// communicating with the SMTP server using the configured credentials and settings.
//func (m *DefaultMailer) sendDirect(message MailMessage) error {
//	addr := fmt.Sprintf("%s:%d", m.smtpConfig.Host, m.smtpConfig.Port)
//
//	// Combine all recipients (To, Cc, Bcc) for the smtp.SendMail `to` parameter.
//	// smtp.SendMail handles de-duplication if necessary and ensures Bcc recipients
//	// are not included in the message headers.
//	var allRecipients []string
//	allRecipients = append(allRecipients, message.To...)
//
//	// Construct the email message content (headers and body).
//	headers := make(map[string]string)
//	headers["From"] = m.smtpConfig.From
//	headers["To"] = strings.Join(message.To, ", ")
//	headers["Subject"] = message.Subject
//	headers["MIME-version"] = "1.0"
//	headers["Content-Type"] = "text/plain; charset=\"UTF-8\"" // Assuming plain text email.
//
//	var emailBuilder strings.Builder
//	for key, value := range headers {
//		emailBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
//	}
//	emailBuilder.WriteString("\r\n") // Standard separator between headers and body.
//	emailBuilder.WriteString(message.Body)
//
//	// Prepare SMTP authentication.
//	auth := smtp.PlainAuth("", m.smtpConfig.Username, m.smtpConfig.Password, m.smtpConfig.Host)
//
//	// Send the email using smtp.SendMail.
//	// Note: smtp.SendMail handles the underlying connection, EHLO, AUTH, MAIL, RCPT, DATA, QUIT sequence.
//	// For more complex scenarios like explicit STARTTLS on a non-standard port, or SMTPS (TLS on connect,
//	// typically port 465), a custom dialer and smtp.NewClient followed by client.StartTLS (if needed)
//	// and client.Auth would be required. This implementation assumes the SMTP server handles TLS
//	// appropriately for the configured port (e.g., implicit TLS for port 465, or STARTTLS upgrade
//	// for port 587 after initial connection).
//	err := smtp.SendMail(addr, auth, m.smtpConfig.From, allRecipients, []byte(emailBuilder.String()))
//	if err != nil {
//		fmt.Printf("failed to send email: %v\n", err)
//		return fmt.Errorf("smtp.SendMail failed for recipients %v: %w", allRecipients, err)
//	}
//	return nil
//}
//
//// Close gracefully shuts down the DefaultMailer.
//// It initiates the shutdown by cancelling the mailer_service's internal context. This signals
//// all worker goroutines to stop processing new tasks and to exit. Close then waits for
//// all workers to complete their current tasks and fully shut down (using sync.WaitGroup).
//// Finally, it closes the task queue, releasing associated resources.
////
//// This method is designed to be idempotent; calling it multiple times will not cause errors
//// or panics, though the actual shutdown operations occur only once.
//// If the parent application context (provided at initialization) is cancelled, the mailer_service
//// will also begin its shutdown sequence automatically. In such cases, calling Close is still
//// recommended (and safe) to ensure a clean wait for full shutdown and resource release.
////
//// After Close returns (or after the parent context is cancelled and shutdown completes),
//// the mailer_service should not be used to send further emails; subsequent calls to Send will fail.
//func (m *DefaultMailer) Close() error {
//	// if m.logger != nil { m.logger.Println("Close called. Cancelling internal context...") }
//	m.cancel() // Signal workers to stop by cancelling the internal context. This is idempotent.
//
//	// if m.logger != nil { m.logger.Println("Waiting for worker goroutines to finish...") }
//	m.wg.Wait() // Wait for all worker goroutines to complete their current tasks and exit.
//	// if m.logger != nil { m.logger.Println("All worker goroutines have finished.") }
//
//	// Close the taskQueue. This must happen after all workers have stopped.
//	// To prevent a panic if Close() is called multiple times (and thus this part is reached multiple times),
//	// we need to ensure close(m.taskQueue) is called only once.
//	// A common way is to use sync.Once, or manage a 'closed' state.
//	// For this implementation, we rely on the fact that after m.wg.Wait(), if the channel
//	// hasn't been closed yet, this is the point to do it.
//	// If Close is called concurrently, m.cancel() is safe. m.wg.Wait() is safe.
//	// The critical part is `close(m.taskQueue)`.
//	// A simple check for idempotency:
//	// We can try to send a non-blocking signal on a private channel to see if it's already closed,
//	// or use a dedicated `sync.Once` for closing the taskQueue.
//	// For now, assuming that `Close` is typically called once by the owner of the mailer_service instance.
//	// If stricter idempotency for channel closing is required, `sync.Once` is the idiomatic Go solution:
//	//   var once sync.Once
//	//   once.Do(func() { close(m.taskQueue) })
//	// However, this would require `once` to be part of the `DefaultMailer` struct.
//	// Given the current structure, if `Close` is called sequentially multiple times,
//	// the second call to `close(m.taskQueue)` would panic.
//	// Let's assume for now that `Close` is called as intended (once, or if multiple times,
//	// the caller handles the fact that only the first call does the full cleanup).
//	// A robust solution would use sync.Once for closing m.taskQueue.
//	// For simplicity in this iteration, we'll proceed without adding sync.Once here,
//	// but acknowledge this as a point for further hardening if strict idempotency of
//	// multiple Close calls is a firm requirement for all internal operations.
//	// The m.cancel() and m.wg.Wait() parts are already idempotent.
//
//	// if m.logger != nil { m.logger.Println("Closing task queue...") }
//	// This will panic if called more than once on a closed channel.
//	// To make it safe for multiple calls, one would typically use sync.Once.
//	// For now, we'll assume Close is called in a way that this is safe (e.g., once).
//	func() {
//		defer func() {
//			// Recover from panic if trying to close an already closed channel.
//			// This makes this specific part of Close idempotent.
//			_ = recover()
//		}()
//		close(m.taskQueue)
//	}()
//	// if m.logger != nil { m.logger.Println("Task queue closed. Mailer shutdown complete.") }
//	return nil
//}
