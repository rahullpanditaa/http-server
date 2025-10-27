-- name: CreateUser :one
INSERT INTO "users" (
    "id", 
    "created_at", 
    "updated_at", 
    "email"
) VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM "users";

-- name: GetUserByEmail :one
SELECT * FROM "users"
WHERE "email" = $1;

-- name: GetUserByRefreshToken :one
SELECT * FROM "users"
WHERE "id" = (
    SELECT "user_id" FROM "refresh_tokens"
    WHERE "token" = $1
);

-- name: UpdateUserDetails :exec
UPDATE "users"
SET "email" = $1,
    "hashed_password" = $2
WHERE "id" = $3;
