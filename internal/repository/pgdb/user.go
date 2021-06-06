package pgdb

import (
    "bot-for-fans/internal/models"
    "errors"
    "github.com/jmoiron/sqlx"
)

type UserRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
    return &UserRepository{db: db}
}

type user struct {
    ID        int64  `db:"id"`
    ChatID    int64  `db:"chat_id"`
    FirstName string `db:"first_name"`
    LastName  string `db:"last_name"`
    UserName  string `db:"user_name"`
}

const getUserByChatIDQuery = `
select *
from telegram.users u
where u.chat_id = $1;
`

func (r *UserRepository) GetByChatID(chatID int64) (models.User, error) {
    var result user
    err := r.db.Get(&result, getUserByChatIDQuery, chatID)
    return models.User{
        ID:        result.ID,
        ChatID:    result.ChatID,
        FirstName: result.FirstName,
        LastName:  result.LastName,
        UserName:  result.UserName,
    }, err
}

const createUserQuery = `
insert into telegram.users (chat_id, first_name, last_name, user_name) 
values ($1, $2, $3, $4) 
returning id;
`

func (r *UserRepository) Create(chatID int64, firstName, lastName, userName string) (id int64, err error) {
    rows, err := r.db.Queryx(createUserQuery, chatID, firstName, lastName, userName)
    defer func() { _ = rows.Close() }()

    if err != nil {
        return 0, err
    }
    if !rows.Next() {
        return 0, errors.New("missing id in returning result")
    }

    err = rows.Scan(&id)
    return
}

const getAllUsersBySubscriptionQuery = `
select 
       u.id,
       u.chat_id,
       u.first_name,
       u.last_name,
       u.user_name
from telegram.users u
left join telegram.user_subscriptions usbscr on u.id = usbscr.user_id
where usbscr.subscription_id = $1;
`

func (r *UserRepository) GetAllBySubscription(subscriptionID int64) ([]models.User, error) {
    var result = make([]user, 0)
    err := r.db.Select(&result, getAllUsersBySubscriptionQuery, subscriptionID)
    if err != nil {
        return []models.User{}, err
    }

    var users = make([]models.User, len(result))
    for i := range result {
        users[i] = models.User{
            ID:        result[i].ID,
            ChatID:    result[i].ChatID,
            FirstName: result[i].FirstName,
            LastName:  result[i].LastName,
            UserName:  result[i].UserName,
        }
    }
    return users, nil
}