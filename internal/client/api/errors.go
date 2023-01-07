package api

import "errors"

var (
	ErrInternalServer    = errors.New("cервер недоступен, попробуйте позднее")
	ErrServerTimout      = errors.New("cервер недоступен, попробуйте позднее")
	ErrSerialization     = errors.New("ошибка сериализации")
	ErrDeserialization   = errors.New("ошибка десериализации")
	ErrRequestPrepare    = errors.New("не удалось подготовить http запрос")
	ErrUserAlreadyExists = errors.New("уже зарегистрирован пользователь с таким логином")
	ErrUserNotFound      = errors.New("пользователь с таким логином не найден")
	ErrAuthRequire       = errors.New("необходим авторизоваться")

	ErrAlreadyExists = errors.New("ID с таким идентификатором уже зарегистрирован")
	ErrNotFound      = errors.New("не найдена запись с таким идентификатором")
)
