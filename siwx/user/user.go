package user

import (
    "github.com/gofiber/fiber/v2"
    "github.com/nnn-community/go-siwx/siwx/db"
    "github.com/nnn-community/go-utils/arrays"
    "strings"
)

type User struct {
    IsLoggedIn  bool                    `json:"IsLoggedIn"`
    ID          string                  `json:"id"`
    Address     string                  `json:"address"`
    Permissions []string                `json:"permissions"`
    UserData    *map[string]interface{} `json:"userData"`
}

var Unauthorized = User{
    IsLoggedIn:  false,
    ID:          "",
    Address:     "",
    Permissions: []string{},
    UserData:    nil,
}

func Get(c *fiber.Ctx) User {
    u := c.Locals("siwx.user").(*User)

    return *u
}

func (u User) Can(permissions []string) bool {
    if len(permissions) == 0 || arrays.Contains(u.Permissions, "is-admin") {
        return true
    }

    _, err := arrays.Find(permissions, func(_ int, value string) bool {
        return arrays.Contains(u.Permissions, value)
    })

    return err == nil
}

func (u User) CanAll(permissions []string) bool {
    if len(permissions) == 0 || arrays.Contains(u.Permissions, "is-admin") {
        return true
    }

    passed := arrays.Filter(permissions, func(_ int, value string) bool {
        return arrays.Contains(u.Permissions, value)
    })

    return len(passed) == len(permissions)
}

type dbUser struct {
    ID      string `json:"id"`
    Address string `json:"address"`
    GroupID string `json:"group_id"`
}

type dbPermission struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

func GetByAddress(address string, getUserData *func(user User) map[string]interface{}) (User, error) {
    var item dbUser

    err := db.Connection.Raw(`
        SELECT u.id, ua.address, u.group_id
		FROM user_addresses ua
		JOIN users u ON ua.user_id = u.id AND u.active = TRUE
		LEFT JOIN "groups" g ON u.group_id = g.id
		WHERE ua.address = ?
        LIMIT 1
    `, strings.ToLower(address)).First(&item).Error

    if err != nil {
        return Unauthorized, err
    }

    var groupPermissions []dbPermission
    var userPermissions []dbPermission
    permissions := []string{}

    db.Connection.Raw(`
        SELECT p.id, p.name
        FROM group_permissions gp
        LEFT JOIN permissions p ON gp.permission_id = p.id
        WHERE gp.group_id = ?
    `, item.GroupID).Scan(&groupPermissions)

    for _, p := range groupPermissions {
        permissions = append(permissions, p.Name)
    }

    db.Connection.Raw(`
        SELECT p.id, p.name
        FROM user_permissions up
        LEFT JOIN permissions p ON up.permission_id = p.id
        WHERE up.user_id = ?
    `, item.ID).Scan(&userPermissions)

    for _, p := range userPermissions {
        permissions = append(permissions, p.Name)
    }

    u := User{
        IsLoggedIn:  true,
        ID:          item.ID,
        Address:     item.Address,
        Permissions: permissions,
        UserData:    nil,
    }

    if getUserData != nil {
        getter := *getUserData
        data := getter(u)

        u.UserData = &data
    }

    return u, nil
}
