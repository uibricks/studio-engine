SET statement_timeout = 60000; -- 60 seconds
SET lock_timeout = 60000; -- 60 seconds

--gopg:split

    CREATE TABLE repositories (
        id serial PRIMARY key,
        created_at timestamptz NULL default current_timestamp,
        deleted_at timestamptz NULL,
        project_id text NULL,
        config jsonb NULL default '{}'::json ,
        project_version int UNIQUE
    );