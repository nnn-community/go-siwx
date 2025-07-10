package db

func Seed() {
    var groupAdminId string
    var groupDefaultId string

    Connection.Raw(`
        INSERT INTO groups (id, name, "default")
        VALUES (gen_random_uuid(), 'Admin', FALSE)
        ON CONFLICT (name) DO UPDATE SET "default" = FALSE
        RETURNING id
    `).Scan(&groupAdminId)

    Connection.Raw(`
        INSERT INTO groups (id, name, "default")
        VALUES (gen_random_uuid(), 'Default', TRUE)
        ON CONFLICT (name) DO UPDATE SET "default" = TRUE
        RETURNING id
    `).Scan(&groupDefaultId)

    var permissionAdminId string

    Connection.Raw(`
        INSERT INTO permissions (id, name)
        VALUES (gen_random_uuid(), 'is-admin')
        ON CONFLICT (name) DO UPDATE SET name = 'is-admin'
        RETURNING id
    `).Scan(&permissionAdminId)

    Connection.Exec(`TRUNCATE group_permissions`)

    Connection.Exec(`
        INSERT INTO group_permissions (id, group_id, permission_id)
        VALUES (gen_random_uuid(), ?, ?)
    `, groupAdminId, permissionAdminId)
}
