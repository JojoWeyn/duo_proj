package service

import (
	"crypto/tls"
	"fmt"
	clientSmtp "github.com/JojoWeyn/duo-proj/identity-service/pkg/client/smtp"
	"math/rand"
	"net/smtp"
	"time"
)

type VerificationService struct {
	client *clientSmtp.SMTPConfig
}

func NewVerificationService(cfg *clientSmtp.SMTPConfig) *VerificationService {
	return &VerificationService{client: cfg}
}

func (vs *VerificationService) GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (vs *VerificationService) SendVerificationCode(email, code string) error {
	tlsConfig := &tls.Config{
		ServerName:         vs.client.Server,
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", vs.client.Server+":"+vs.client.Port, tlsConfig)
	if err != nil {
		return fmt.Errorf("error connecting to SMTP server: %v", err)
	}

	client, err := smtp.NewClient(conn, vs.client.Server)
	if err != nil {
		conn.Close()
		return fmt.Errorf("error creating SMTP client: %v", err)
	}

	auth := smtp.PlainAuth("", vs.client.Sender, vs.client.Password, vs.client.Server)
	if err := client.Auth(auth); err != nil {
		client.Close()
		return fmt.Errorf("authentication failed: %v", err)
	}

	if err := client.Mail(vs.client.Sender); err != nil {
		client.Close()
		return fmt.Errorf("failed to set sender: %v", err)
	}

	if err := client.Rcpt(email); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %v", err)
	}
	defer wc.Close()

	body := fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <title>Email Confirmation</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                background-color: #f4f4f9;
                margin: 0;
                padding: 20px;
            }
            .email-container {
                background-color: white;
                border-radius: 8px;
                padding: 20px;
                box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                max-width: 600px;
                margin: auto;
            }
            h1 {
                color: #333;
                font-size: 24px;
            }
            p {
                color: #555;
                font-size: 16px;
            }
            .code {
                font-size: 24px;
                font-weight: bold;
                color:rgb(63, 63, 63);
                background-color:rgb(255, 175, 77);
                padding: 10px;
                border-radius: 4px;
            }
            .footer {
                font-size: 12px;
                color: #999;
                text-align: center;
                margin-top: 20px;
            }
        </style>
    </head>
    <body>

        <div class="email-container">
            <h1>Код проверки электронной почты</h1>
            <p>Здравствуйте,</p>
            <p>Введите этот код на экране проверки личности:</p>

            <p class="code">%s</p>

            <p>Срок действия этого кода истекает в ближайшее время. <br> 
			Если вы не можете найти экран проверки личности, попробуйте войти в систему еще раз. <br>
			Если вы не пытались войти в свою учетную запись, мы рекомендуем вам сбросить пароль прямо сейчас.</p>

            <div class="footer">
                <p>Kozhura</p>
            </div>
        </div>

    </body>
    </html>
    `, code)

	msg := "From: " + vs.client.Sender + "\n" +
		"To: " + email + "\n" +
		"Subject: " + "Kozhura Код проверки" + "\n" +
		"MIME-Version: 1.0" + "\n" +
		"Content-Type: text/html; charset=UTF-8" + "\n\n" +
		body

	_, err = wc.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write email content: %v", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close data writer: %v", err)
	}

	return nil
}
