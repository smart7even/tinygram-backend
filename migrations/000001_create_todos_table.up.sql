CREATE TABLE todos (
  id BIGSERIAL PRIMARY KEY,
  name varchar(255) NOT NULL,
  complete BOOLEAN NOT NULL,
)