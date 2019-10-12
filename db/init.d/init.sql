CREATE DATABASE test;
GRANT ALL PRIVILEGES ON DATABASE test TO root;
CREATE TABLE users (
  id varchar(255) PRIMARY KEY,
  username varchar(255) UNIQUE,
  imsValue varchar(255),
  imsType varchar(255)
);
CREATE TABLE emails (
  user_id varchar(255),
  type varchar(255),
  value varchar(255)
);

INSERT INTO users (id, username, imsValue, imsType)
VALUES ('0001', 'andy@example.com', 'andy@foo.com', 'xmpp');
INSERT INTO emails (user_id, type, value)
VALUES ('0001', 'work', 'andy@example.com'), ('0001', 'home', 'andy@aol.com');

INSERT INTO users (id, username, imsValue, imsType)
VALUES ('0002', 'someone@example.org', 'ifc', 'home');
INSERT INTO emails (user_id, type, value)
VALUES ('0002', 'work', 'someone@example.org'), ('0002', 'home', 'someone@hotmail.com');

INSERT INTO users (id, username, imsValue, imsType)
VALUES ('0003', 'someone@fox.com', 'fox', 'fun');
INSERT INTO emails (user_id, type, value)
VALUES ('0003', 'work', 'someone@fox.com'), ('0003', 'home', 'someone@abc.com');

INSERT INTO users (id, username, imsValue, imsType)
VALUES ('0004', 'a@foo.com', 'irc', 'booooo');
INSERT INTO emails (user_id, type, value)
VALUES ('0004', 'work', 'a@foo.com'), ('0004', 'home', 'b@abc.com');
