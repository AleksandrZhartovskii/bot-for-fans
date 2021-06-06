package schemas

type Response struct {
    Error bool        `json:"error"`
    Data  interface{} `json:"data"`
}

type ErrorResponse struct {
    Message string `json:"message"`
}