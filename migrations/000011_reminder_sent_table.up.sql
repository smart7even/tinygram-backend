CREATE TABLE reminder_sent (
  id BIGSERIAL PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  reminder_id BIGINT NOT NULL,
  device_id BIGINT NOT NULL,
  sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users (id),
  FOREIGN KEY (reminder_id) REFERENCES reminders (id),
  FOREIGN KEY (device_id) REFERENCES devices (id)
);
