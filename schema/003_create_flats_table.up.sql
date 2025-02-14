CREATE TYPE flat_status AS ENUM ('created', 'approved', 'declined', 'on moderation');

CREATE TABLE IF NOT EXISTS flats (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    number INT NOT NULL,
    house_id BIGINT NOT NULL,
    price INT NOT NULL,
    rooms INT NOT NULL,
    status flat_status NOT NULL
);
