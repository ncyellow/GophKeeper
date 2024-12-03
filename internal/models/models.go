// Package models contains structures that describe our domain entities
// right now all data is strings, we should replace them with correct types + proper validation
package models

// User - user type
type User struct {
	UserID   int64  `json:"-"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Card - bank card
type Card struct {
	UserID   int64  `json:"-"`
	ID       string `json:"id"`
	FIO      string `json:"fio"` // The name on the card may differ from the actual name
	Number   string `json:"number"`
	Date     string `json:"date"`
	CVV      string `json:"cvv"`
	MetaInfo string `json:"metainfo"`
}

// Text - text content
type Text struct {
	UserID   int64  `json:"-"`
	ID       string `json:"id"`
	Content  string `json:"content"`
	MetaInfo string `json:"metainfo"`
}

// Binary - binary data
type Binary struct {
	UserID   int64  `json:"-"`
	ID       string `json:"id"`
	Data     []byte `json:"data"`
	MetaInfo string `json:"metainfo"`
}

// Login - login data
type Login struct {
	UserID   int64  `json:"-"`
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	MetaInfo string `json:"metainfo"`
}
