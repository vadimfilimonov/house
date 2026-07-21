CREATE TABLE IF NOT EXISTS house_subscriptions (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    house_id BIGINT NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_house_subscriptions_house FOREIGN KEY (house_id) REFERENCES houses(id),
    CONSTRAINT unique_house_subscription UNIQUE (house_id, email)
);
