package models

type User struct {
    ID        int64
    ChatID    int64
    FirstName string
    LastName  string
    UserName  string
}

func (u User) GetAppeal() string {
    if u.FirstName != "" && u.LastName != "" {
        return u.FirstName + " " + u.LastName
    }
    if u.FirstName != "" {
        return u.FirstName
    }
    if u.LastName != "" {
        return u.LastName
    }
    if u.UserName != "" {
        return u.UserName
    }
    return "Dear Fan"
}