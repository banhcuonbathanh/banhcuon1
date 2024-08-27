-- Migration to update users table

-- Drop existing table if it exists


CREATE TABLE reading_tests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    test_number INTEGER NOT NULL,
    sections JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);