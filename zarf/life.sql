CREATE USER wooknight;

GRANT ALL PRIVILEGES ON DATABASE life TO wooknight;
CREATE TABLE life_user (
    id serial primary key,
    username varchar(255) unique not null ,
    password VARCHAR ( 50 ) NOT NULL,
    email varchar(255) unique not null ,
    created_on TIMESTAMP not null ,
    updated_on TIMESTAMP not null 
);

CREATE TABLE goal (
    id serial primary key,
    goal varchar(255) unique not null ,
    created_on TIMESTAMP not null ,
    updated_on TIMESTAMP not null 
);

CREATE TABLE thought (
    id serial primary key,
    thought text not null ,
    created_on TIMESTAMP not null ,
    updated_on TIMESTAMP not null 
);

CREATE TABLE habit (
    id serial primary key,
    habit varchar(255) unique not null ,
    created_on TIMESTAMP not null ,
    updated_on TIMESTAMP not null 
);

CREATE TABLE tag (
    id serial primary key,
    tagName varchar(255) unique not null ,
    created_on TIMESTAMP not null ,
    updated_on TIMESTAMP not null 
);


CREATE TABLE model (
    id serial primary key,
    modelName varchar(255) unique not null ,
    steps text[] ,
    created_on TIMESTAMP not null ,
    updated_on TIMESTAMP not null 
);