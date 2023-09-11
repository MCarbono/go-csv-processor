drop table if exists movies; 

create table movies (
    id integer primary key, 
    title text, 
    year text,
    genres text
);
