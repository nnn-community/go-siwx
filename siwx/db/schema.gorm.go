package db

import (
    "time"
)

type Permission struct {
    ID               string `gorm:"type:varchar(36);primaryKey;default:uuid_generate_v4()"`
    Name             string `gorm:"type:varchar(32);unique"`
    GroupPermissions []GroupPermission
    UserPermissions  []UserPermission
    CreatedAt        time.Time
    UpdatedAt        time.Time
}

func (Permission) TableName() string {
    return "permissions"
}

type Group struct {
    ID               string `gorm:"type:varchar(36);primaryKey;default:uuid_generate_v4()"`
    Name             string `gorm:"unique"`
    GroupPermissions []GroupPermission
    Users            []User
    IsDefault        bool `gorm:"default:false;column:default"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
}

func (Group) TableName() string {
    return "groups"
}

type GroupPermission struct {
    ID           string `gorm:"type:varchar(36);primaryKey;default:uuid_generate_v4()"`
    GroupID      string `gorm:"type:varchar(36);uniqueIndex:group_permission_identifier"`
    Group        Group
    PermissionID string `gorm:"type:varchar(36);uniqueIndex:group_permission_identifier"`
    Permission   Permission
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

func (GroupPermission) TableName() string {
    return "group_permissions"
}

type UserPermission struct {
    ID           string `gorm:"type:varchar(36);primaryKey;default:uuid_generate_v4()"`
    UserID       string `gorm:"type:varchar(36);uniqueIndex:user_permission_identifier"`
    User         User
    PermissionID string `gorm:"type:varchar(36);uniqueIndex:user_permission_identifier"`
    Permission   Permission
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

func (UserPermission) TableName() string {
    return "user_permissions"
}

type User struct {
    ID              string `gorm:"type:varchar(36);primaryKey;default:uuid_generate_v4()"`
    Addresses       []UserAddress
    UserPermissions []UserPermission
    GroupID         string `gorm:"type:varchar(36)"`
    Group           Group
    Active          bool `gorm:"default:true"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

func (User) TableName() string {
    return "users"
}

type UserAddress struct {
    ID        string `gorm:"type:varchar(36);primaryKey;default:uuid_generate_v4()"`
    Address   string `gorm:"type:varchar(42);unique"`
    UserID    string `gorm:"type:varchar(36)"`
    User      User
    Master    bool `gorm:"default:false"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (UserAddress) TableName() string {
    return "user_addresses"
}
