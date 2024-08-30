-- Create enum for QuestionType
CREATE TYPE question_type AS ENUM ('MultipleChoice', 'TrueFalseNotGiven', 'Matching', 'ShortAnswer');

-- Create reading_req_models table
CREATE TABLE reading_req_models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reading_req_test_type JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create reading_res_models table
CREATE TABLE reading_res_models (
    id BIGSERIAL PRIMARY KEY,
    reading_res_type JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

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

-- Create page_request_models table
CREATE TABLE page_request_models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    page_number INTEGER NOT NULL,
    page_size INTEGER NOT NULL
);