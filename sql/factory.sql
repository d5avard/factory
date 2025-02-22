CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    question VARCHAR(100) NOT NULL
);

INSERT INTO questions (question)
VALUES 
('My first question');