-- +goose Up
-- SQL in section 'Up' is applied when this migration is applied

CREATE TABLE IF NOT EXISTS userinfo (
    id bigserial PRIMARY KEY,
    username text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    token text NOT NULL,
    login_time timestamp DEFAULT now(),
    logout_time timestamp,
    is_active boolean 
);

CREATE TABLE IF NOT EXISTS chatroom (
    id bigserial PRIMARY KEY,
    room_name text NOT NULL
);

CREATE TABLE IF NOT EXISTS chatroom_user (
    chatroom_id bigint REFERENCES chatroom(id),
    user_id bigint REFERENCES userinfo(id),
);

CREATE TABLE IF NOT EXISTS message (
    id bigserial PRIMARY KEY,
    chatroom_id bigint REFERENCES chatroom(id),
    user_id bigint REFERENCES userinfo(id),
    message text NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
);

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back
DROP TABLE message;
DROP TABLE chatroom_user;
DROP TABLE chatroom;
DROP TABLE user;

