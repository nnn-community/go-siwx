package siwx

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/gofiber/storage/redis/v3"
    "net/url"
    "strings"
)

type sessionData struct {
    SiwxUser *User `json:"siwx-user"`
}

var unauthorizedUser = User{
    IsLoggedIn:  false,
    ID:          "",
    Address:     "",
    Permissions: []string{},
    UserData:    nil,
}

func getSessionUser(r *redis.Storage, sessionId string) (User, error) {
    decoded, err := url.QueryUnescape(sessionId)

    if err != nil {
        return unauthorizedUser, err
    }

    parts := strings.Split(decoded, ".")

    if len(parts) < 1 {
        return unauthorizedUser, errors.New("invalid cookie format")
    }

    s, err := r.Get("sess:" + parts[0])

    if err != nil {
        return unauthorizedUser, err
    }

    if string(s) == "" {
        return unauthorizedUser, fmt.Errorf("empty JSON string")
    }

    var data sessionData

    if err := json.Unmarshal(s, &data); err != nil {
        return unauthorizedUser, fmt.Errorf("invalid JSON: %w", err)
    }

    if data.SiwxUser == nil {
        return unauthorizedUser, fmt.Errorf("siwxUser field missing")
    }

    user := *data.SiwxUser
    user.IsLoggedIn = true

    return user, nil
}
