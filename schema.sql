CREATE TABLE todos (
    id bigint GENERATED ALWAYS AS IDENTITY,
    name VARCHAR NOT NULL,
	complete boolean not NULL
);