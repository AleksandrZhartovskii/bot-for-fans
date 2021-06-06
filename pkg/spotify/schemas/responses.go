package schemas

type ErrorResp struct {
    Error struct {
        Status  int32  `json:"status"`
        Message string `json:"message"`
    } `json:"error"`
}

type ErrorAuth struct {
    Error string `json:"error"`
}

type GetArtistsByNameResp struct {
    Artists struct{
        Items []Artist `json:"items"`
    } `json:"artists"`
}

type GetArtistsAlbumsResp struct {
    Albums []Album `json:"items"`
}

type GetSeveralAlbumsTracksResp struct {
    Albums []struct{
        Tracks struct{
            Items []Track `json:"items"`
        } `json:"tracks"`
    } `json:"albums"`
}

type AuthResponse struct {
    AccessToken string `json:"access_token"`
}