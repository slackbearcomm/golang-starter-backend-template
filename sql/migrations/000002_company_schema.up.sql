BEGIN;

-- Organizations
CREATE TABLE "organizations" (
    "id" bigserial UNIQUE NOT NULL PRIMARY KEY,
    "uid" uuid UNIQUE NOT NULL,
    "code" varchar UNIQUE NOT NULL,
    "name" varchar NOT NULL,
    "website" varchar,
    "logo" JSONB,
    "sector" varchar NOT NULL,
    "status" varchar NOT NULL DEFAULT '',
    "is_final" boolean NOT NULL DEFAULT FALSE,
    "is_archived" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamptz NOT NULL DEFAULT NOW(),
    "updated_at" timestamptz NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON organizations
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TABLE "departments" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "code" varchar UNIQUE NOT NULL,
    "org_uid" uuid NOT NULL REFERENCES organizations (uid),
    "name" varchar NOT NULL,
    "status" varchar NOT NULL DEFAULT '',
    "is_final" boolean NOT NULL DEFAULT FALSE,
    "is_archived" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamptz NOT NULL DEFAULT NOW(),
    "updated_at" timestamptz NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON departments
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TABLE "roles" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "code" varchar UNIQUE NOT NULL,
    "org_uid" uuid NOT NULL REFERENCES organizations (uid),
    "department_id" bigint NOT NULL REFERENCES departments (id),
    "name" varchar NOT NULL,
    "permissions" text[],
    "is_management" boolean NOT NULL DEFAULT FALSE,
    "status" varchar NOT NULL DEFAULT '',
    "is_final" boolean NOT NULL DEFAULT FALSE,
    "is_archived" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamptz NOT NULL DEFAULT NOW(),
    "updated_at" timestamptz NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON roles
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Users
CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "first_name" varchar NOT NULL,
    "last_name" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "phone" varchar UNIQUE NOT NULL,
    "is_admin" boolean NOT NULL DEFAULT FALSE,
    "org_uid" uuid REFERENCES organizations (uid),
    "role_id" bigint REFERENCES roles (id),
    "status" varchar NOT NULL DEFAULT '',
    "is_final" boolean NOT NULL DEFAULT FALSE,
    "is_archived" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamptz NOT NULL DEFAULT NOW(),
    "updated_at" timestamptz NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TABLE "otp_sessions" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "user_id" bigint NOT NULL REFERENCES users (id),
    "token" varchar NOT NULL,
    "is_valid" boolean NOT NULL DEFAULT FALSE,
    "expires_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT NOW(),
    "updated_at" timestamptz NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON otp_sessions
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TABLE "auth_sessions" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "user_id" bigint NOT NULL REFERENCES users (id),
    "token" uuid UNIQUE NOT NULL,
    "is_valid" boolean NOT NULL DEFAULT FALSE,
    "expires_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT NOW(),
    "updated_at" timestamptz NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON auth_sessions
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Activities
CREATE TABLE "user_activities" (
    "id" bigserial PRIMARY KEY NOT NULL,
    "user_id" bigint NOT NULL REFERENCES users (id),
    "org_uid" uuid REFERENCES organizations (uid),
    "action" varchar NOT NULL,
    "object_id" bigint,
    "object_type" varchar,
    "session_token" uuid NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT NOW(),
    "updated_at" timestamptz NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON user_activities
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

COMMIT;