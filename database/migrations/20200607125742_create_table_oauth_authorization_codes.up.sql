create table oauth_authorization_codes (
    code varchar(100) not null,
    user_id integer not null references users(id),
    expires_at timestamp with time zone not null
);

create unique index unique_active_oauth_authorization_codes on oauth_authorization_codes(code, user_id);
