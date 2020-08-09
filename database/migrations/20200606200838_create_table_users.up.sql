create table users (
    id serial4 primary key,
    name varchar(300) not null,
    email varchar(500) not null,
    google_login_email varchar(500),
    github_login_username varchar(200),
    created_at timestamp with time zone not null default 'now()',
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

create unique index unique_email_on_users on users(email);
create unique index unique_google_login_email_on_users on users(google_login_email);
create unique index unique_github_login_username_on_users on users(github_login_username);

create type workspace_access_level as enum ('owner', 'editor', 'viewer');

create table workspace_users (
    id serial4 primary key,
    workspace_id integer not null references workspaces(id),
    user_id integer not null references users(id),
    access_level workspace_access_level not null,
    created_at timestamp with time zone not null default 'now()',
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

create unique index unique_user_on_workspace on workspace_users(workspace_id, user_id) where deleted_at is not null;

create type project_access_level as enum ('editor', 'viewer');

create table project_users (
    id serial4 primary key,
    project_id integer not null references projects(id),
    user_id integer not null references users(id),
    access_level project_access_level not null,
    created_at timestamp with time zone not null default 'now()',
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

create unique index unique_user_on_project on project_users(project_id, user_id) where deleted_at is not null;
