create type resource_state as enum ('unknown', 'provisioning', 'provisioning_failed', 'provisioning_success');

create table resources (
    id serial4 primary key,
    name varchar(200) not null,
    type varchar(100) not null,
    payload jsonb not null,
    state resource_state not null default 'unknown',
    created_at timestamp with time zone not null default 'now()',
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

create table resource_logs (
    resource_id bigint not null references resources(id),
    payload varchar(2000) not null,
    created_at timestamp with time zone not null default 'now()'
);

create index order_resource_logs on resource_logs(resource_id, created_at desc);
