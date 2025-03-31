package siwx

import "github.com/nnn-community/go-utils/arrays"

type User struct {
    IsLoggedIn  bool                    `json:"isLoggedIn"`
    ID          string                  `json:"id"`
    Address     string                  `json:"address"`
    Permissions []string                `json:"permissions"`
    UserData    *map[string]interface{} `json:"userData"`
}

func (u User) Can(required []string) bool {
    if len(required) == 0 || arrays.Contains(u.Permissions, "is-admin") {
        return true
    }

    _, err := arrays.Find(required, func(_ int, value string) bool {
        return arrays.Contains(u.Permissions, value)
    })

    return err == nil
}

func (u User) CanAll(required []string) bool {
    if len(required) == 0 || arrays.Contains(u.Permissions, "is-admin") {
        return true
    }

    passed := arrays.Filter(required, func(_ int, value string) bool {
        return arrays.Contains(u.Permissions, value)
    })

    return len(passed) == len(required)
}
