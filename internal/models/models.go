package models

// User - тип пользователя
type User struct {
	UserID   int64  `json:"-"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Card - карта банка
type Card struct {
	UserID   int    `json:"-"`
	ID       string `json:"id"`
	FIO      string `json:"fio"` // Написание ФИО на карте может отличаться от реального
	Number   string `json:"number"`
	Date     string `json:"date"`
	CVV      string `json:"cvv"`
	MetaInfo string `json:"metainfo"`
}

// Text - текстовый контент
type Text struct {
	UserID   int    `json:"-"`
	ID       string `json:"id"`
	Content  string `json:"content"`
	MetaInfo string `json:"metainfo"`
}

// Binary - бинарные данные
type Binary struct {
	UserID   int    `json:"-"`
	ID       string `json:"id"`
	Data     []byte `json:"data"`
	MetaInfo string `json:"metainfo"`
}

// Login - данные по логинам
type Login struct {
	UserID   int64  `json:"-"`
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	MetaInfo string `json:"metainfo"`
}
