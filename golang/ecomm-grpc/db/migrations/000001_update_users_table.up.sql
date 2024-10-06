CREATE TYPE question_type AS ENUM ('MultipleChoice', 'TrueFalseNotGiven', 'Matching', 'ShortAnswer');

-- Recreate tables
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'Employee',
    phone VARCHAR(20),
    image VARCHAR(255),
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_role ON users(role);

ALTER TABLE users
ADD CONSTRAINT check_valid_role
CHECK (role IN ('Admin', 'Employee', 'Manager', 'Guest'));

CREATE TABLE IF NOT EXISTS sessions (
    id varchar(255) PRIMARY KEY NOT NULL,
    user_email VARCHAR(100) NOT NULL,
    refresh_token TEXT NOT NULL,
    is_revoked BOOLEAN DEFAULT false,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS reading_test_models (
    id SERIAL PRIMARY KEY,
    test_number INTEGER NOT NULL,
    sections JSONB NOT NULL
);

CREATE TABLE IF NOT EXISTS section_models (
    id SERIAL PRIMARY KEY,
    section_number INTEGER NOT NULL,
    time_allowed INTEGER NOT NULL,
    passages JSONB NOT NULL
);

CREATE TABLE IF NOT EXISTS passage_models (
    id SERIAL PRIMARY KEY,
    passage_number INTEGER NOT NULL,
    title TEXT NOT NULL,
    content JSONB NOT NULL,
    questions JSONB NOT NULL
);

CREATE TABLE IF NOT EXISTS paragraph_content_models (
    id SERIAL PRIMARY KEY,
    paragraph_summary TEXT NOT NULL,
    key_words TEXT NOT NULL,
    key_sentence TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS question_models (
    id SERIAL PRIMARY KEY,
    question_number INTEGER NOT NULL,
    type question_type NOT NULL,
    content TEXT NOT NULL,
    options TEXT[] NULL,
    correct_answer JSONB NULL
);

CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    author_id BIGINT NOT NULL,
    parent_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    avatar VARCHAR(255),
    role VARCHAR(50) DEFAULT 'Employee',
    owner_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS dishes (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    description TEXT NOT NULL,
    image VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'Available',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS dish_snapshots (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    description TEXT NOT NULL,
    image VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'Available',
    dish_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS tables (
    number INTEGER PRIMARY KEY,
    capacity INTEGER NOT NULL,
    status VARCHAR(50) DEFAULT 'Available',
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- CREATE TABLE IF NOT EXISTS tables (
--     number INTEGER PRIMARY KEY,
--     capacity INTEGER NOT NULL,
--     status INTEGER DEFAULT 0, -- 0 corresponds to AVAILABLE in the enum
--     token VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
-- );

CREATE TABLE IF NOT EXISTS guests (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    table_number INTEGER,
    refresh_token VARCHAR(255),
    refresh_token_expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (table_number) REFERENCES tables(number) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    guest_id BIGINT,
    table_number INTEGER,
    dish_snapshot_id BIGINT UNIQUE NOT NULL,
    quantity INTEGER NOT NULL,
    order_handler_id BIGINT,
    status VARCHAR(50) DEFAULT 'Pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (guest_id) REFERENCES guests(id) ON DELETE SET NULL,
    FOREIGN KEY (table_number) REFERENCES tables(number) ON DELETE SET NULL,
    FOREIGN KEY (dish_snapshot_id) REFERENCES dish_snapshots(id) ON DELETE CASCADE,
    FOREIGN KEY (order_handler_id) REFERENCES accounts(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    token VARCHAR(255) PRIMARY KEY,
    account_id BIGINT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sockets (
    socket_id VARCHAR(255) PRIMARY KEY,
    account_id BIGINT UNIQUE,
    guest_id BIGINT UNIQUE,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE SET NULL,
    FOREIGN KEY (guest_id) REFERENCES guests(id) ON DELETE SET NULL
);

ALTER TABLE accounts
ADD CONSTRAINT fk_owner
FOREIGN KEY (owner_id) REFERENCES accounts(id) ON DELETE SET NULL;