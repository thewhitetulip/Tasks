CREATE TABLE task (
id integer primary key autoincrement,
title varchar(100),
content text,
is_deleted char(1) default 'N',
created_date timestamp,
last_modified_at timestamp,
finish_date timestamp
, priority int);

CREATE TABLE files(name varchar(1000) not null, autoName varchar(255) not null);
