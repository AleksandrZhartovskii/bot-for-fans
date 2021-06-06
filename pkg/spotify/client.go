package spotify

import (
    "bot-for-fans/pkg/spotify/schemas"
    "encoding/base64"
    "errors"
    "fmt"
    "github.com/mailru/easyjson"
    "github.com/valyala/fasthttp"
    "strings"
    "time"
)

type Client struct {
    secret              string
    clientID            string
    auth                string
    getArtistsURI       string
    getArtistsAlbumsURI string
    getSeveralAlbumsURI string
}

type Config struct {
    AppSecret           string
    AppClientID         string
    GetArtistsURI       string
    GetArtistsAlbumsURI string
    GetSeveralAlbumsURI string
}

func NewClient(c Config) *Client {
    return &Client{
        secret:              c.AppSecret,
        clientID:            c.AppClientID,
        getArtistsURI:       c.GetArtistsURI,
        getArtistsAlbumsURI: c.GetArtistsAlbumsURI,
        getSeveralAlbumsURI: c.GetSeveralAlbumsURI,
    }
}

func (c *Client) GetArtistsByName(name string) ([]schemas.Artist, error) {
    if err := c.refreshToken(); err != nil {
        return []schemas.Artist{}, err
    }

    name = strings.Trim(name, " ")
    name = strings.ReplaceAll(name, " ", "%20")
    uri := fmt.Sprintf("%s?q=%s&type=artist&limit=5", c.getArtistsURI, name)

    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)

    req.Header.SetMethod(fasthttp.MethodGet)
    req.Header.Set(fasthttp.HeaderAuthorization, c.auth)
    req.SetRequestURI(uri)

    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    if err := fasthttp.DoTimeout(req, resp, time.Second*5); err != nil {
        return []schemas.Artist{}, err
    }

    if resp.StatusCode() != fasthttp.StatusOK {
        var response schemas.ErrorResp
        if err := easyjson.Unmarshal(resp.Body(), &response); err != nil {
            return []schemas.Artist{}, err
        }
        return []schemas.Artist{}, errors.New(response.Error.Message)
    }

    var response schemas.GetArtistsByNameResp
    if err := easyjson.Unmarshal(resp.Body(), &response); err != nil {
        return []schemas.Artist{}, err
    }
    return response.Artists.Items, nil
}

func (c *Client) GetArtistsAlbums(artistID string) ([]schemas.Album, error) {
    if err := c.refreshToken(); err != nil {
        return []schemas.Album{}, err
    }

    uri := fmt.Sprintf("%s/%s/albums", c.auth, artistID)

    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)

    req.Header.SetMethod(fasthttp.MethodGet)
    req.Header.Set(fasthttp.HeaderAuthorization, c.auth)
    req.SetRequestURI(uri)

    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    if err := fasthttp.DoTimeout(req, resp, time.Second*5); err != nil {
        return []schemas.Album{}, err
    }

    if resp.StatusCode() != fasthttp.StatusOK {
        var response schemas.ErrorResp
        if err := easyjson.Unmarshal(resp.Body(), &response); err != nil {
            return []schemas.Album{}, err
        }
        return []schemas.Album{}, errors.New(response.Error.Message)
    }

    var response schemas.GetArtistsAlbumsResp
    if err := easyjson.Unmarshal(resp.Body(), &response); err != nil {
        return []schemas.Album{}, err
    }
    return response.Albums, nil
}

func (c *Client) GetSeveralAlbumsTracks(albumsIDS ...string) ([]schemas.Track, error) {
    if err := c.refreshToken(); err != nil {
        return []schemas.Track{}, err
    }

    var ids string
    for i, id := range albumsIDS {
        ids += id
        if i != len(albumsIDS)-1 {
            ids += ","
        }
    }
    uri := fmt.Sprintf("%s?ids=%s", c.auth, ids)

    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)

    req.Header.SetMethod(fasthttp.MethodGet)
    req.Header.Set(fasthttp.HeaderAuthorization, c.auth)
    req.SetRequestURI(uri)

    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    if err := fasthttp.DoTimeout(req, resp, time.Second*5); err != nil {
        return []schemas.Track{}, err
    }

    if resp.StatusCode() != fasthttp.StatusOK {
        var response schemas.ErrorResp
        if err := easyjson.Unmarshal(resp.Body(), &response); err != nil {
            return []schemas.Track{}, err
        }
        return []schemas.Track{}, errors.New(response.Error.Message)
    }

    var response schemas.GetSeveralAlbumsTracksResp
    if err := easyjson.Unmarshal(resp.Body(), &response); err != nil {
        return []schemas.Track{}, err
    }

    var tracks []schemas.Track
    for _, album := range response.Albums {
        tracks = append(tracks, album.Tracks.Items...)
    }
    return tracks, nil
}

func (c *Client) refreshToken() error {
    auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.clientID, c.secret)))

    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)

    req.Header.SetMethod(fasthttp.MethodPost)
    req.SetRequestURI("https://accounts.spotify.com/api/token")
    req.Header.Set(fasthttp.HeaderAuthorization, "Basic "+auth)
    req.SetBodyString("grant_type=client_credentials")

    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    if err := fasthttp.DoTimeout(req, resp, time.Second*5); err != nil {
        return err
    }

    if resp.StatusCode() != fasthttp.StatusOK {
        var response schemas.ErrorAuth
        if err := easyjson.Unmarshal(resp.Body(), &response); err != nil {
            return err
        }
        return errors.New(response.Error)
    }

    var response schemas.AuthResponse
    if err := easyjson.Unmarshal(resp.Body(), &response); err != nil {
        return err
    }

    c.auth = "Bearer " + response.AccessToken
    return nil
}