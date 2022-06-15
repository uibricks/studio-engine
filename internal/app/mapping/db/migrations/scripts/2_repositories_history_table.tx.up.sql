SET statement_timeout = 60000; -- 60 seconds
SET lock_timeout = 60000; -- 60 seconds

--gopg:split

CREATE TABLE repositories_history (
	id serial PRIMARY key,
	project_id text NULL,
	config jsonb NULL,
	project_version int NULL,
	repos_created_at timestamptz NULL,
	repos_deleted_at timestamptz NULL,
	created_at timestamptz NULL default current_timestamp
);