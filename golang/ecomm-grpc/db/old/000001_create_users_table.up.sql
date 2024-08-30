CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TYPE question_type AS ENUM ('MultipleChoice', 'TrueFalseNotGiven', 'Matching', 'ShortAnswer');


-- Create reading_test_models table
CREATE TABLE reading_test_models (
    id SERIAL PRIMARY KEY,
    test_number INTEGER NOT NULL,
    sections JSONB NOT NULL
);

-- Create section_models table
CREATE TABLE section_models (
    id SERIAL PRIMARY KEY,
    section_number INTEGER NOT NULL,
    time_allowed INTEGER NOT NULL,
    passages JSONB NOT NULL
);

-- Create passage_models table
CREATE TABLE passage_models (
    id SERIAL PRIMARY KEY,
    passage_number INTEGER NOT NULL,
    title TEXT NOT NULL,
    content JSONB NOT NULL,
    questions JSONB NOT NULL
);

-- Create paragraph_content_models table
CREATE TABLE paragraph_content_models (
    id SERIAL PRIMARY KEY,
    paragraph_summary TEXT NOT NULL,
    key_words TEXT NOT NULL,
    key_sentence TEXT NOT NULL
);

-- Create question_models table
CREATE TABLE question_models (
    id SERIAL PRIMARY KEY,
    question_number INTEGER NOT NULL,
    type question_type NOT NULL,
    content TEXT NOT NULL,
    options TEXT[] NULL,
    correct_answer JSONB NULL
);

