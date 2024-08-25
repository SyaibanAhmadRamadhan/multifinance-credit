-- Create the payment_methods table
CREATE TABLE payment_methods
(
    id        INT AUTO_INCREMENT PRIMARY KEY,
    name      VARCHAR(255) NOT NULL,
    code_name VARCHAR(255) NOT NULL,
    enabled   BOOLEAN
);