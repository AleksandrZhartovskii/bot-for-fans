package tracker

import (
    "bot-for-fans/internal/models"
    botapi "bot-for-fans/internal/server/schemas"
    "bot-for-fans/pkg/spotify/schemas"
    "errors"
    "github.com/valyala/fasthttp"
    "go.uber.org/zap"
    "time"
)

func (t *T) checkSubscriptions() {
    t.logger.Info("Subscriptions checking successfully started")

    subscriptions, err := t.repo.GetAll()
    if err != nil {
        t.logger.Error(
            "Error getting all subscriptions",
            zap.String("error", err.Error()),
        )
        return
    }

    for _, s := range subscriptions {
        if err := t.checkSubscription(s); err != nil {
            t.logger.Error(
                "Error checking subscription",
                zap.String("error", err.Error()),
                zap.Int64("subscription_id", s.ID),
                zap.String("author_name", s.AuthorName),
                zap.String("author_spotify_id", s.AuthorSpotifyID),
                zap.String("last_update_time", s.LastUpdateTime.String()),
            )
        }
    }

    t.logger.Info("Subscriptions checking successfully completed")
}

func (t *T) checkSubscription(s models.Subscription) error {
    albums, err := t.spotify.GetArtistsAlbums(s.AuthorSpotifyID)
    if err != nil {
        return err
    }

    rawAlbumsIDS, err := t.getRawAlbumsIDS(s, albums)
    if err != nil {
        return err
    }
    if len(rawAlbumsIDS) == 0 {
        return t.repo.ChangeLastUpdateTime(s.ID, time.Now())
    }

    newTracks, err := t.spotify.GetSeveralAlbumsTracks(rawAlbumsIDS...)
    if err != nil {
        return err
    }

    for _, track := range newTracks {
        if err := t.sendToBot(s.ID, track.ExternalURLS.Spotify); err != nil {
            return err
        }
    }
    return nil
}

func (t *T) getRawAlbumsIDS(s models.Subscription, a []schemas.Album) ([]string, error) {
    var rawAlbums []string
    for _, album := range a {
        switch album.ReleaseDatePrecision {
        case schemas.DatePrecisionYear:
            // These date format are not processed as it is impossible
            // to accurately determine the newness of the album
            continue
        case schemas.DatePrecisionMonth:
            // These date format are not processed as it is impossible
            // to accurately determine the newness of the album
            continue
        case schemas.DatePrecisionDay:
            albumDate, err := time.Parse("2006-01-02", album.ReleaseDate)
            if err != nil {
                return rawAlbums, err
            }

            sdatestr := s.LastUpdateTime.Format("2006-01-02")
            lastSubscriptionUpdateDate, err := time.Parse("2006-01-02", sdatestr)
            if err != nil {
                return rawAlbums, err
            }

            if albumDate.After(lastSubscriptionUpdateDate) {
                rawAlbums = append(rawAlbums, album.ID)
            }

        default:
            return rawAlbums, errors.New("unknown date precisions")
        }
    }
    return rawAlbums, nil
}

func (t *T) sendToBot(subcrID int64, musicURL string) error {
   var body = botapi.Notification{
        SubscriptionID: subcrID,
        MusicURL:       musicURL,
    }
    b, err := body.MarshalJSON()
    if err != nil {
        return err
    }

    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)

    req.Header.SetMethod(fasthttp.MethodPost)
    req.SetRequestURI(t.botURL)
    req.SetBody(b)

    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    return fasthttp.DoTimeout(req, resp, time.Second*5)
}