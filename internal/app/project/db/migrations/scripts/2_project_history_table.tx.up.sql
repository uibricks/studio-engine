SET statement_timeout = 60000; -- 60 seconds
SET lock_timeout = 60000; -- 60 seconds

--gopg:split

CREATE TABLE project_history (
     id serial PRIMARY key,
     luid uuid not null,
     name varchar(200),
     client_id varchar(200),
     version int,
     user_version varchar(200),
     config jsonb NULL,
     state varchar(100),
     created_by int,
     updated_by int,
     project_created_at timestamptz NULL,
     project_updated_at timestamptz NULL,
     project_deleted_at timestamptz NULL,
     created_at timestamptz NULL default current_timestamp
);
