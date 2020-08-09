create table workspace_repository_connections (
    id serial4 primary key,
    workspace_id integer not null references workspaces(id),
    identifier varchar(500) not null,
    provider varchar(200) not null,
    payload jsonb not null,
    created_at timestamp with time zone not null default 'now()',
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

create unique index unique_repository_provider_per_workspace on workspace_repository_connections(workspace_id, provider, identifier);
