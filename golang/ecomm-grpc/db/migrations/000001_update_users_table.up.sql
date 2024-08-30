CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_admin BOOLEAN DEFAULT false,
    phone BIGINT,
    image VARCHAR(255),
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
    id varchar(255) PRIMARY KEY NOT NULL,
    user_email VARCHAR(100) NOT NULL,
    refresh_token TEXT NOT NULL,
    is_revoked BOOLEAN DEFAULT false,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);
CREATE TYPE question_type AS ENUM ('MultipleChoice', 'TrueFalseNotGiven', 'Matching', 'ShortAnswer');


-- Create reading_test_models table
CREATE TABLE IF NOT EXISTS reading_test_models (
    id SERIAL PRIMARY KEY,
    test_number INTEGER NOT NULL,
    sections JSONB NOT NULL
);

-- Create section_models table
CREATE TABLE IF NOT EXISTS section_models (
    id SERIAL PRIMARY KEY,
    section_number INTEGER NOT NULL,
    time_allowed INTEGER NOT NULL,
    passages JSONB NOT NULL
);

-- Create passage_models table
CREATE TABLE IF NOT EXISTS passage_models (
    id SERIAL PRIMARY KEY,
    passage_number INTEGER NOT NULL,
    title TEXT NOT NULL,
    content JSONB NOT NULL,
    questions JSONB NOT NULL
);

-- Create paragraph_content_models table
CREATE TABLE IF NOT EXISTS paragraph_content_models (
    id SERIAL PRIMARY KEY,
    paragraph_summary TEXT NOT NULL,
    key_words TEXT NOT NULL,
    key_sentence TEXT NOT NULL
);

-- Create question_models table
CREATE TABLE IF NOT EXISTS question_models (
    id SERIAL PRIMARY KEY,
    question_number INTEGER NOT NULL,
    type question_type NOT NULL,
    content TEXT NOT NULL,
    options TEXT[] NULL,
    correct_answer JSONB NULL
);

