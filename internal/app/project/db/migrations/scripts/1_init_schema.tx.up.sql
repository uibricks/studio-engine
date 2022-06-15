SET statement_timeout = 60000; -- 60 seconds
SET lock_timeout = 60000; -- 60 seconds

--gopg:split

CREATE TABLE objects (
     id serial PRIMARY key,
     name varchar(200),
     type int,
     parent_id int,
     client_id varchar(200),
     created_at timestamptz NULL default current_timestamp,
     updated_at timestamptz NULL default current_timestamp,
     deleted_at timestamptz NULL,
     created_by int,
     updated_by int,
     state varchar(100),
     version serial,
     luid uuid not null,
     user_version varchar(200),
     config jsonb NULL,
     constraint fk_object foreign key(parent_id) references objects(id)
);




