CREATE TABLE expressions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    expression TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    result FLOAT DEFAULT NULL
);