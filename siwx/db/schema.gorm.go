package db

import (
    "time"
)

type Permission struct {
    ID               string `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid();index:,type:HASH"`
    Name             string `gorm:"type:varchar(32);unique"`
    GroupPermissions []GroupPermission
    UserPermissions  []UserPermission
    CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Permission) TableName() string {
    return "permissions"
}

type Group struct {
    ID               string `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid();index:,type:HASH"`
    Name             string `gorm:"unique"`
    GroupPermissions []GroupPermission
    Users            []User
    IsDefault        bool      `gorm:"default:false;column:default;index:,type:HASH"`
    CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Group) TableName() string {
    return "groups"
}

type GroupPermission struct {
    ID           string `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid();index:,type:HASH"`
    GroupID      string `gorm:"type:varchar(36);column:group_id;uniqueIndex:group_permission_identifier;index:,type:HASH"`
    Group        Group
    PermissionID string `gorm:"type:varchar(36);column:permission_id;uniqueIndex:group_permission_identifier;index:,type:HASH"`
    Permission   Permission
    CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (GroupPermission) TableName() string {
    return "group_permissions"
}

type UserPermission struct {
    ID           string `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid();index:,type:HASH"`
    UserID       string `gorm:"type:varchar(36);column:user_id;uniqueIndex:user_permission_identifier;index:,type:HASH"`
    User         User
    PermissionID string `gorm:"type:varchar(36);column:permission_id;uniqueIndex:user_permission_identifier;index:,type:HASH"`
    Permission   Permission
    CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserPermission) TableName() string {
    return "user_permissions"
}

type User struct {
    ID              string `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid();index:,type:HASH"`
    Addresses       []UserAddress
    UserPermissions []UserPermission
    GroupID         string `gorm:"type:varchar(36);column:group_id;index:,type:HASH"`
    Group           Group
    Active          bool      `gorm:"default:true;index:,type:HASH"`
    CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (User) TableName() string {
    return "users"
}

type UserAddress struct {
    ID        string `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid();index:,type:HASH"`
    Address   string `gorm:"type:varchar(42);unique;index:,type:HASH"`
    UserID    string `gorm:"type:varchar(36);column:user_id;index:,type:HASH"`
    User      User
    Master    bool      `gorm:"default:false;index:,type:HASH"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserAddress) TableName() string {
    return "user_addresses"
}
