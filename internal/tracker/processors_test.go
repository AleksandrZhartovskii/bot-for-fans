package tracker

import (
    "fmt"
    "log"
    "testing"
    "time"
)

func Test_ParseDay(t *testing.T) {
    date := "1991-12-01"
    datetime, err := time.Parse("2006-01-02", date)
    if err != nil {
        log.Fatal(err)
    }

    nowtime, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(datetime)
    fmt.Println(nowtime)
}