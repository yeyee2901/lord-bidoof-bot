package datasource

import (
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/yeyee2901/lord-bidoof-bot/pkg/config"
)

type DataSource struct {
	Config *config.AppConfig
	DB     *sqlx.DB
	Redis  *redis.Client
}

func NewDataSource(c *config.AppConfig, db *sqlx.DB, r *redis.Client) *DataSource {
	return &DataSource{c, db, r}
}

func (ds *DataSource) InsertPrivateChatToDB(chat *PrivateChat) error {
	q := `
        INSERT INTO telegram_private_chat
            (chat_id, name, username, bio)
        VALUES
            (:chat_id, :name, :username, :bio)
    `

	tx, err := ds.DB.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(q, chat)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (ds *DataSource) GetPrivateChat(chatId int64) (*PrivateChat, error) {
	var args []any
	args = append(args, chatId)

	q := `
        SELECT
            name
        FROM
            telegram_private_chat
        WHERE
            chat_id = ?
    `

	res := new(PrivateChat)
	err := ds.DB.Get(res, q, args...)

	return res, err
}
