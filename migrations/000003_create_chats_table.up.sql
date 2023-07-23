CREATE TABLE chats (
    id char(36) PRIMARY KEY,
    name char(255) NOT NULL
);

CREATE TABLE chat_user (
    chat_id char(36) NOT NULL,
    user_id char(36) NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chats (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE messages (
    id char(36) PRIMARY KEY,
    chat_id char(36) NOT NULL,
    user_id char(36) NOT NULL,
    text char(255) NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chats (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);