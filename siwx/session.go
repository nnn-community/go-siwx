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
    Siwx     *Siwx `json:"siwx"`
    SiwxUser *User `json:"siwx-user"`
}

var unauthorizedSiwx = Siwx{
    Address: "",
    ChainId: 0,
}

var unauthorizedUser = User{
    IsLoggedIn:  false,
    ID:          "",
    Address:     "",
    Permissions: []string{},
    UserData:    nil,
}

var unauthorizedSession = sessionData{
    Siwx:     &unauthorizedSiwx,
    SiwxUser: &unauthorizedUser,
}

func getSessionSiwx(r *redis.Storage, sessionId string) (Siwx, error) {
    data, err := getSession(r, sessionId)

    return *data.Siwx, err
}

func getSessionUser(r *redis.Storage, sessionId string) (User, error) {
    data, err := getSession(r, sessionId)
    user := *data.SiwxUser

    if err == nil {
        user.IsLoggedIn = true
    }

    return user, err
}

func parseSessionId(sessionId string) (*string, error) {
    decoded, err := url.QueryUnescape(sessionId)

    if err != nil {
        return nil, err
    }

    parts := strings.Split(decoded, ".")

    if len(parts) < 1 {
        return nil, errors.New("invalid cookie format")
    }

    id := "sess:" + parts[0]

    return &id, nil
}

func getSession(r *redis.Storage, sessionId string) (sessionData, error) {
    sessId, err := parseSessionId(sessionId)

    if err != nil {
        return unauthorizedSession, err
    }

    s, err := r.Get(*sessId)

    if err != nil {
        return unauthorizedSession, err
    }

    if string(s) == "" {
        return unauthorizedSession, fmt.Errorf("empty JSON string")
    }

    var data sessionData

    if err := json.Unmarshal(s, &data); err != nil {
        return unauthorizedSession, fmt.Errorf("invalid JSON: %w", err)
    }

    if data.SiwxUser == nil {
        return unauthorizedSession, fmt.Errorf("siwxUser field missing")
    }

    return data, nil
}
