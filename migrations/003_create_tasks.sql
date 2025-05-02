CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    expression_id INT NOT NULL REFERENCES expressions(id),
    arg1 FLOAT NOT NULL,
    arg2 FLOAT NOT NULL,
    operation VARCHAR(10) NOT NULL,
    operation_time INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending'
);