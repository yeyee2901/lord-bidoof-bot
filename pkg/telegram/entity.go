package telegram

type RespGetMe struct {
	Id                      uint64
	IsBot                   bool  
	FirstName               string
	LastName                string
	Username                string
	CanJoinGroups           bool  
	CanReadAllGroupMessages bool  
}

type RespSendMessage struct {
	MessageID int64
	Recipient string
}
