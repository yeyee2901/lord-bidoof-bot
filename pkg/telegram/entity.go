package telegram

type RespGetMe struct {
	Ok          bool         `json:"ok" yaml:"ok"`
	Description string       `json:"description" yaml:"description"`
	Result      TelegramUser `json:"result" yaml:"result"`
}

type TelegramUser struct {
	Id                      uint64 `json:"id" yaml:"id"`
	IsBot                   bool   `json:"is_bot" yaml:"is_bot"`
	FirstName               string `json:"first_name" yaml:"first_name"`
	LastName                string `json:"last_name" yaml:"last_name"`
	Username                string `json:"username" yaml:"username"`
	CanJoinGroups           bool   `json:"can_join_groups" yaml:"can_join_groups"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages" yaml:"can_read_all_group_messages"`
}
