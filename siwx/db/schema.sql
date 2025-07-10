-- Function to update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ language 'plpgsql';

-- Permissions Table
CREATE TABLE "permissions" (
    "id" VARCHAR(36) PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" VARCHAR(32) UNIQUE NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX ON "permissions" USING HASH ("id");
CREATE TRIGGER update_permissions_updated_at BEFORE UPDATE ON "permissions" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- Groups Table
CREATE TABLE "groups" (
    "id" VARCHAR(36) PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" VARCHAR(255) UNIQUE NOT NULL,
    "default" BOOLEAN NOT NULL DEFAULT false,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX ON "groups" USING HASH ("id");
CREATE TRIGGER update_groups_updated_at BEFORE UPDATE ON "groups" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- Users Table
CREATE TABLE "users" (
    "id" VARCHAR(36) PRIMARY KEY DEFAULT uuid_generate_v4(),
    "group_id" VARCHAR(36) NOT NULL,
    "active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY ("group_id") REFERENCES "groups"("id") ON DELETE RESTRICT ON UPDATE CASCADE
);
CREATE INDEX ON "users" USING HASH ("id");
CREATE INDEX ON "users" USING HASH ("group_id");
CREATE INDEX ON "users" USING HASH ("active");
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON "users" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- Group Permissions Table
CREATE TABLE "group_permissions" (
    "id" VARCHAR(36) PRIMARY KEY DEFAULT uuid_generate_v4(),
    "group_id" VARCHAR(36) NOT NULL,
    "permission_id" VARCHAR(36) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE ("group_id", "permission_id"),
    FOREIGN KEY ("group_id") REFERENCES "groups"("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY ("permission_id") REFERENCES "permissions"("id") ON DELETE RESTRICT ON UPDATE CASCADE
);
CREATE INDEX ON "group_permissions" USING HASH ("id");
CREATE INDEX ON "group_permissions" USING HASH ("group_id");
CREATE INDEX ON "group_permissions" USING HASH ("permission_id");
CREATE TRIGGER update_group_permissions_updated_at BEFORE UPDATE ON "group_permissions" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- User Permissions Table
CREATE TABLE "user_permissions" (
    "id" VARCHAR(36) PRIMARY KEY DEFAULT uuid_generate_v4(),
    "user_id" VARCHAR(36) NOT NULL,
    "permission_id" VARCHAR(36) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE ("user_id", "permission_id"),
    FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY ("permission_id") REFERENCES "permissions"("id") ON DELETE RESTRICT ON UPDATE CASCADE
);
CREATE INDEX ON "user_permissions" USING HASH ("id");
CREATE INDEX ON "user_permissions" USING HASH ("user_id");
CREATE INDEX ON "user_permissions" USING HASH ("permission_id");
CREATE TRIGGER update_user_permissions_updated_at BEFORE UPDATE ON "user_permissions" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- User Addresses Table
CREATE TABLE "user_addresses" (
    "id" VARCHAR(36) PRIMARY KEY DEFAULT uuid_generate_v4(),
    "address" VARCHAR(42) UNIQUE NOT NULL,
    "user_id" VARCHAR(36) NOT NULL,
    "master" BOOLEAN NOT NULL DEFAULT false,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE
);
CREATE INDEX ON "user_addresses" USING HASH ("id");
CREATE INDEX ON "user_addresses" USING HASH ("address");
CREATE INDEX ON "user_addresses" USING HASH ("user_id");
CREATE INDEX ON "user_addresses" USING HASH ("master");
CREATE TRIGGER update_user_addresses_updated_at BEFORE UPDATE ON "user_addresses" FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();