create table oauth_tokens (
    id serial4 not null,
    access_token_hash bytea not null,
    refresh_token_hash bytea not null,
    user_id integer not null references users(id),
    expires_at timestamp with time zone not null,
    requester_ip inet not null,
    requester_user_agent varchar(2000) not null,
    deleted_at timestamp with time zone
);

create unique index unique_active_oauth_tokens on oauth_tokens(access_token_hash, refresh_token_hash, user_id);
