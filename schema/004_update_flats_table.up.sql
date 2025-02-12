ALTER TABLE flats
    ADD CONSTRAINT unique_number_house_id UNIQUE (number, house_id),
    ADD CONSTRAINT fk_house FOREIGN KEY (house_id) REFERENCES houses(id);
