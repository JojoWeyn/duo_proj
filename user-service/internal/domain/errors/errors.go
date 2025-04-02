package errors

import "errors"

var (
	// Ошибки, связанные с пользователем (1000-1999)
	ErrUserNotFound                = errors.New("1001") // Пользователь не найден
	ErrInvalidVerificationCode     = errors.New("1002") // Неверный код подтверждения
	ErrEmailAlreadyExists          = errors.New("1003") // Email уже существует
	ErrInvalidEmailOrPassword      = errors.New("1004") // Неверный email или пароль
	ErrEmailNotConfirmed           = errors.New("1005") // Email не подтвержден
	ErrInvalidRefreshToken         = errors.New("1006") // Неверный токен обновления
	ErrEmailEmpty                  = errors.New("1007") // Email не может быть пустым
	ErrEmailTooLong                = errors.New("1008") // Email слишком длинный
	ErrInvalidEmailFormat          = errors.New("1009") // Неверный формат email
	ErrPasswordTooShort            = errors.New("1010") // Пароль должен содержать не менее 8 символов
	ErrPasswordTooLong             = errors.New("1011") // Пароль слишком длинный
	ErrPasswordNoUppercase         = errors.New("1012") // Пароль должен содержать хотя бы одну заглавную букву
	ErrPasswordNoLowercase         = errors.New("1013") // Пароль должен содержать хотя бы одну строчную букву
	ErrPasswordNoNumber            = errors.New("1014") // Пароль должен содержать хотя бы одну цифру
	ErrPasswordNoSpecialCharacter  = errors.New("1015") // Пароль должен содержать хотя бы один специальный символ
	ErrFailedToGetUserAchievements = errors.New("1016") // Не удалось получить достижения пользователя
	ErrFailedToParseCondition      = errors.New("1017") // Не удалось распарсить условие
	ErrFailedToGetAchievements     = errors.New("1018") // Не удалось получить достижения
	ErrFailedToGetLeaderboard      = errors.New("1019") // Не удалось получить лидерборд
	ErrProgressNotFound            = errors.New("1020") // Прогресс не найден
	ErrFailedToGetUsers            = errors.New("1021") // Не удалось получить пользователей
	ErrFailedToOpenAvatarFile      = errors.New("1022") // Не удалось открыть файл аватара
	ErrFailedToUpdateAvatar        = errors.New("1023") // Не удалось обновить аватар

	// Ошибки, связанные с токенами (2000-2999)
	ErrEmptyUserID             = errors.New("2001") // userID не может быть пустым
	ErrInvalidToken            = errors.New("2002") // Неверный токен
	ErrUnexpectedSigningMethod = errors.New("2003") // Неожиданный метод подписания токена
	ErrUUIDNotFoundInToken     = errors.New("2004") // UUID не найден в токене
	ErrInvalidUUIDFormat       = errors.New("2005") // Неверный формат UUID

	// Ошибки, связанные с SMTP (3000-3999)
	ErrConnectionToSMTPServer    = errors.New("3001") // Ошибка подключения к SMTP серверу
	ErrCreatingSMTPClient        = errors.New("3002") // Ошибка при создании SMTP клиента
	ErrAuthenticationFailed      = errors.New("3003") // Ошибка аутентификации
	ErrFailedToSetSender         = errors.New("3004") // Не удалось установить отправителя
	ErrFailedToSetRecipient      = errors.New("3005") // Не удалось установить получателя
	ErrFailedToGetDataWriter     = errors.New("3006") // Не удалось получить writer для данных
	ErrFailedToWriteEmailContent = errors.New("3007") // Не удалось записать содержимое письма
	ErrFailedToCloseDataWriter   = errors.New("3008") // Не удалось закрыть writer для данных

	// HTTP ошибки (4000-4999)
	ErrNoTokenProvided              = errors.New("4001") // Не предоставлен токен
	ErrEmailRequired                = errors.New("4002") // Email обязателен
	ErrFailedToSaveVerificationCode = errors.New("4003") // Не удалось сохранить код подтверждения
	ErrVerificationCodeInvalid      = errors.New("4004") // Неверный код подтверждения
	ErrFailedToSendVerificationCode = errors.New("4005") // Не удалось отправить код подтверждения
	ErrPasswordResetFailed          = errors.New("4006") // Не удалось сбросить пароль
	ErrInternalServerError          = errors.New("4007") // Внутренняя ошибка сервера

	// HTTP сообщения об успехе (5000-5999)
	SuccessUserRegistered            = "5001: user registered successfully" // Пользователь зарегистрирован успешно
	SuccessPasswordReset             = "5002: password reset successfully"  // Пароль успешно сброшен
	SuccessEmailConfirmed            = "5003: email confirmed successfully" // Email успешно подтвержден
	SuccessVerificationCodeGenerated = "5004: verification code generated"  // Код подтверждения сгенерирован

	// Ошибки, связанные с внешними сервисами (6000-6999)
	ErrServiceUnavailable      = errors.New("6001") // Сервис недоступен
	ErrServiceTimeout          = errors.New("6002") // Время ожидания сервиса истекло
	ErrInvalidServiceResponse  = errors.New("6003") // Неверный ответ от сервиса
	ErrServiceConnectionFailed = errors.New("6004") // Ошибка подключения к внешнему сервису
	ErrExternalServiceError    = errors.New("6005") // Ошибка внешнего сервиса
	ErrInvalidServiceURL       = errors.New("6006") // Неверный URL внешнего сервиса
	ErrTooManyRequests         = errors.New("6007") // Слишком много запросов

	// Ошибки, связанные с валидацией данных (7000-7999)
	ErrInvalidDataFormat     = errors.New("7001") // Неверный формат данных
	ErrDataParsingError      = errors.New("7002") // Ошибка парсинга данных
	ErrMissingRequiredFields = errors.New("7003") // Отсутствуют обязательные поля
	ErrInvalidFieldValue     = errors.New("7004") // Неверное значение поля
	ErrInvalidLimit          = errors.New("7005") // Неверное значение limit
	ErrInvalidOffset         = errors.New("7006") // Неверное значение offset

)
