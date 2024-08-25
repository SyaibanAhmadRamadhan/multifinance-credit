CREATE TABLE limits
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    consumer_id INT NOT NULL,
    tenor       INT NOT NULL,
    amount DOUBLE NOT NULL,
    remaining_amount DOUBLE NOT NULL,

    FOREIGN KEY (consumer_id) REFERENCES consumers (id) ON DELETE CASCADE
);

CREATE INDEX index_limits_consumer_id ON limits (consumer_id);