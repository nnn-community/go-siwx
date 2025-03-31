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
    SiwxUser *User `json:"siwxUser"`
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

    parts := strings.SplitN(decoded, ".", 2)

    if len(parts) != 2 {
        return unauthorizedUser, errors.New("invalid cookie format")
    }

    exp := strings.SplitN(parts[0], ":", 2)

    if len(parts) != 2 {
        return unauthorizedUser, errors.New("invalid cookie format")
    }

    s, err := r.Get("sess:" + exp[1])

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
