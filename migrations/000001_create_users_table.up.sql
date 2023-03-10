CREATE TABLE users (
    username varchar(255) not null unique primary key ,
    bio text,
    last_online timestamp,
    password_hash varchar(512) not null,
    public_key varchar(256) not null unique
);

CREATE TABLE tokens (
    id serial primary key,
    user_username  varchar(255) references users(username),
    token varchar(128) not null unique,
    expires_at date not null
);

CREATE TABLE conversations (
    id serial primary key,
    title varchar(255) not null
);

CREATE TABLE conversation_members (
    conv_id int references conversations(id),
    user_username varchar(255) references users(username),
    unique (conv_id, user_username)
);

CREATE TABLE messages (
    id serial primary key,
    id_in_conv int,
    sender_username varchar(255) references users(username),
    conv_id int references conversations(id),
    send_date timestamp
);

CREATE TABLE message_text (
  id int references messages(id),
  text varchar(5000) not null,
  for_user varchar(255) references users(username)
);
