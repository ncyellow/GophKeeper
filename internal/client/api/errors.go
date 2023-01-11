package api

import "errors"

// Перечень ошибок которые может генерировать client при взаимодействии с сервером
// При это часть ошибок при работе через http клиент приходит как StatusCode и ошибки нет.
// Тогда мы генерируем ошибку, если же grpc прилетает ошибка мы ее прокидываем используя формат ошибки
var (
	FmtErrInternalServer    = "cервер недоступен, попробуйте позднее: %w"
	FmtErrServerTimout      = "cервер недоступен, попробуйте позднее: %w"
	FmtErrDeserialization   = "ошибка десериализации: %w"
	FmtErrRequestPrepare    = "не удалось подготовить http запрос: %w"
	FmtErrUserAlreadyExists = "уже зарегистрирован пользователь с таким логином: %w"
	FmtErrUserNotFound      = "пользователь с таким логином не найден: %w"
	FmtErrAlreadyExists     = "ID с таким идентификатором уже зарегистрирован: %w"
	FmtErrNotFound          = "не найдена запись с таким идентификатором: %w"
	FmtErrSerialization     = "ошибка сериализации: %w"

	ErrSerialization     = errors.New("ошибка сериализации")
	ErrAuthRequire       = errors.New("необходимо авторизоваться")
	ErrUserAlreadyExists = errors.New("уже зарегистрирован пользователь с таким логином")
	ErrInternalServer    = errors.New("cервер недоступен, попробуйте позднее")
	ErrUserNotFound      = errors.New("пользователь с таким логином не найден")
	ErrAlreadyExists     = errors.New("ID с таким идентификатором уже зарегистрирован")
	ErrNotFound          = errors.New("не найдена запись с таким идентификатором")
)
