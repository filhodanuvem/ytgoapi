CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
create table posts (
    id uuid primary key default uuid_generate_v4(),
    username varchar(255) not null,
    body text not null,
    created_at timestamp not null default now()
);