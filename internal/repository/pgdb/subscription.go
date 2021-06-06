package pgdb

import (
    "bot-for-fans/internal/models"
    "errors"
    "github.com/jmoiron/sqlx"
    "time"
)

type SubscriptionRepository struct {
    db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) *SubscriptionRepository {
    return &SubscriptionRepository{db: db}
}

const getOrCreateSubscriptionQuery = `
with neworexisting as (
    insert into telegram.subscriptions (author_name, author_spotify_id) values ($1, $2)
    on conflict (author_spotify_id) do update set exist = true
    returning id
)
select id from neworexisting;
`

func (r *SubscriptionRepository) GetOrCreate(authorName, authorID string) (id int64, err  error) {
    rows, err := r.db.Queryx(getOrCreateSubscriptionQuery, authorName, authorID)
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

type subscription struct {
    ID              int64     `db:"id"`
    AuthorName      string    `db:"author_name"`
    AuthorSpotifyID string    `db:"author_spotify_id"`
    LastUpdateTime  time.Time `db:"last_update_time"`
}

const getAllSubscriptionsQuery = `
select * from telegram.subscriptions;
`

func (r *SubscriptionRepository) GetAll() ([]models.Subscription, error) {
    var result = make([]subscription, 0)
    err := r.db.Select(&result, getAllSubscriptionsQuery)
    if err != nil {
        return []models.Subscription{}, err
    }

    var subscriptions = make([]models.Subscription, len(result))
    for i := range result {
        subscriptions[i] = models.Subscription{
            ID:              result[i].ID,
            AuthorSpotifyID: result[i].AuthorSpotifyID,
            AuthorName:      result[i].AuthorName,
            LastUpdateTime:  result[i].LastUpdateTime,
        }
    }
    return subscriptions, nil
}

const changeLastUpdateTimeSubscriptionQuery = `
update telegram.subscriptions set last_update_time = $1 where id = $2;
`

func (r *SubscriptionRepository) ChangeLastUpdateTime(id int64, t time.Time) error {
    rows, err := r.db.Queryx(changeLastUpdateTimeSubscriptionQuery, t, id)
    defer func() { _ = rows.Close() }()
    return err
}