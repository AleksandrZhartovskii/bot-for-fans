package schemas

type Artist struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type DatePrecision string
const (
    DatePrecisionYear  DatePrecision = "year"  // example: 1981
    DatePrecisionMonth DatePrecision = "month" // example: 1981-12
    DatePrecisionDay   DatePrecision = "day"   // example: 1981-12-15
)

type Album struct {
    ID                   string        `json:"id"`
    ReleaseDate          string        `json:"release_date"`
    ReleaseDatePrecision DatePrecision `json:"release_date_precision"`
}

type Track struct {
    ExternalURLS struct {
        Spotify string `json:"spotify"`
    } `json:"external_urls"`
}