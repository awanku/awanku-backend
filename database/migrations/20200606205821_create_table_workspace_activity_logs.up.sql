create table workspace_activity_logs (
    workspace_id integer not null references workspaces(id),
    payload varchar(1000) not null,
    created_at timestamp with time zone not null default 'now()'
);

create index order_on_workspace_activity_logs on workspace_activity_logs(workspace_id, created_at desc);
