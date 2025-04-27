package service_test

import (
	"fmt"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/service"
	clientSmtp "github.com/JojoWeyn/duo-proj/identity-service/pkg/client/smtp"
	"github.com/stretchr/testify/require"
	"sync"

	"net"
	"net/smtp"

	"testing"
	"time"
)

func setupMockSMTPServer(t *testing.T) (string, string, net.Listener, func()) {
	listener, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err, "failed to start mock SMTP server")

	addr := listener.Addr().String()
	host, port, err := net.SplitHostPort(addr)
	require.NoError(t, err, "failed to split host and port")

	return host, port, listener, func() { listener.Close() }
}

// --- GenerateVerificationCode ---

func TestGenerateVerificationCode_Success(t *testing.T) {
	cfg := &clientSmtp.SMTPConfig{}
	vs := service.NewVerificationService(cfg)

	code := vs.GenerateVerificationCode()
	require.Len(t, code, 6, "code should be 6 digits")
	_, err := fmt.Sscanf(code, "%d", new(int))
	require.NoError(t, err, "code should be numeric")

	// Проверяем, что коды разные
	time.Sleep(1 * time.Millisecond)
	code2 := vs.GenerateVerificationCode()
	require.NotEqual(t, code, code2, "codes should be different")
}

// --- SendVerificationCode ---

func TestSendVerificationCode_Success(t *testing.T) {
	host, port, listener, cleanup := setupMockSMTPServer(t)
	defer cleanup()

	cfg := &clientSmtp.SMTPConfig{
		Server:   host,
		Port:     port,
		Sender:   "test@kozhura.com",
		Password: "testpass",
	}
	vs := service.NewVerificationService(cfg)

	// Используем WaitGroup для синхронизации
	var wg sync.WaitGroup
	wg.Add(1)

	// Запускаем мок-сервер
	go func() {
		defer wg.Done()
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, cfg.Server)
		if err != nil {
			return
		}
		defer client.Close()

		err = client.Auth(smtp.PlainAuth("", cfg.Sender, cfg.Password, cfg.Server))
		if err != nil {
			return
		}
		err = client.Mail(cfg.Sender)
		if err != nil {
			return
		}
		err = client.Rcpt("recipient@example.com")
		if err != nil {
			return
		}
		wc, err := client.Data()
		if err != nil {
			return
		}
		wc.Close()
	}()

	// Даем серверу время запуститься
	time.Sleep(100 * time.Millisecond)

	err := vs.SendVerificationCode("recipient@example.com", "123456")
	require.NoError(t, err, "should send email successfully")

	// Ждем завершения работы мока
	wg.Wait()
}

func TestSendVerificationCode_InvalidEmail(t *testing.T) {
	host, port, listener, cleanup := setupMockSMTPServer(t)
	defer cleanup()

	cfg := &clientSmtp.SMTPConfig{
		Server:   host,
		Port:     port,
		Sender:   "test@kozhura.com",
		Password: "testpass",
	}
	vs := service.NewVerificationService(cfg)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, cfg.Server)
		if err != nil {
			return
		}
		defer client.Close()

		err = client.Auth(smtp.PlainAuth("", cfg.Sender, cfg.Password, cfg.Server))
		if err != nil {
			return
		}
		err = client.Mail(cfg.Sender)
		if err != nil {
			return
		}
		// Симулируем ошибку для невалидного email
		err = client.Rcpt("invalid-email")
		if err == nil {
			client.Quit()
			return
		}
	}()

	time.Sleep(100 * time.Millisecond)

	err := vs.SendVerificationCode("invalid-email", "123456")
	require.Error(t, err, "should fail with invalid email")
	require.Contains(t, err.Error(), "failed to set recipient", "should mention recipient error")

	wg.Wait()
}

func TestSendVerificationCode_ConnectionFailure(t *testing.T) {
	cfg := &clientSmtp.SMTPConfig{
		Server:   "invalid-server",
		Port:     "9999",
		Sender:   "test@kozhura.com",
		Password: "testpass",
	}
	vs := service.NewVerificationService(cfg)

	err := vs.SendVerificationCode("recipient@example.com", "123456")
	require.Error(t, err, "should fail to connect")
	require.Contains(t, err.Error(), "error connecting to SMTP server", "should mention connection error")
}

func TestSendVerificationCode_AuthFailure(t *testing.T) {
	host, port, listener, cleanup := setupMockSMTPServer(t)
	defer cleanup()

	cfg := &clientSmtp.SMTPConfig{
		Server:   host,
		Port:     port,
		Sender:   "test@kozhura.com",
		Password: "wrongpass",
	}
	vs := service.NewVerificationService(cfg)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, cfg.Server)
		if err != nil {
			return
		}
		defer client.Close()

		// Симулируем ошибку авторизации
		err = client.Auth(smtp.PlainAuth("", cfg.Sender, "invalid-password", cfg.Server))
		if err == nil {
			client.Quit()
			return
		}
	}()

	time.Sleep(100 * time.Millisecond)

	err := vs.SendVerificationCode("recipient@example.com", "123456")
	require.Error(t, err, "should fail authentication")
	require.Contains(t, err.Error(), "authentication failed", "should mention authentication error")

	wg.Wait()
}
