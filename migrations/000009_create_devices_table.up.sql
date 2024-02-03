CREATE TABLE devices (
    id BIGSERIAL PRIMARY KEY,
    device_id VARCHAR(100) NOT NULL,
    device_token VARCHAR(100) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    device_os VARCHAR(100) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);