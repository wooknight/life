CREATE USER wooknight;

GRANT ALL PRIVILEGES ON DATABASE life TO wooknight;
CREATE TABLE users (
    id serial primary key,
    user_name varchar(255) unique not null ,
    password VARCHAR ( 50 ) NOT NULL,
    email varchar(255) unique not null ,
    created_at TIMESTAMP not null ,
    updated_at TIMESTAMP not null 
);

CREATE TABLE goals (
    id serial primary key,
    goal_name varchar(255) unique not null ,
    created_at TIMESTAMP not null ,
    updated_at TIMESTAMP not null 
);

CREATE TABLE thoughts (
    id serial primary key,
    thought_name text not null ,
    created_at TIMESTAMP not null ,
    updated_at TIMESTAMP not null 
);

CREATE TABLE habits (
    id serial primary key,
    habit_name varchar(255) unique not null ,
    habit_desc text not null ,
    created_at TIMESTAMP not null ,
    updated_at TIMESTAMP not null 
);

CREATE TABLE tags (
    id serial primary key,
    tag_name varchar(255) unique not null ,
    created_at TIMESTAMP not null ,
    updated_at TIMESTAMP not null 
);


CREATE TABLE models (
    id serial primary key,
    model_name varchar(255) unique not null ,
    steps text[] ,
    created_at TIMESTAMP not null ,
    updated_at TIMESTAMP not null 
);