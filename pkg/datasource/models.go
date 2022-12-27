package datasource

type PrivateChat struct {
	ChatID   int64  `json:"chat_id" db:"chat_id"`
	Username string `json:"username" db:"username"`
	Name     string `json:"name" db:"name"`
	Bio      string `json:"bio" db:"bio"`
}
