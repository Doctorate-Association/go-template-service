package models

type User struct {

	// ID uuid
	ID string `json:"ID"`

	// Name user name
	Name string `json:"name"`

	// DisplayName user display name
	DisplayName string `json:"displayName"`

	// Email user email
	Email string `json:"email"`

	// Avatar user avatar
	Avatar string `json:"avatar"`

	// PHone user phone
	Phone string `json:"phone"`

	// Exp user expire time
	Exp float64 `json:"exp"`

	// Nvb user not valid beforeß
	Nvb float64 `json:"nvb"`

	// Iat  user issued at
	Iat float64 `json:"iat"`

	// Jti user jwt id
	Jti string `json:"jti"`
}
