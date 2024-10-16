-- Drop existing tables and types
DROP TABLE IF EXISTS set_snapshot_dishes;
DROP TABLE IF EXISTS set_snapshots;
DROP TABLE IF EXISTS set_dishes;
DROP TABLE IF EXISTS sets;
DROP TABLE IF EXISTS sockets;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS guests;
DROP TABLE IF EXISTS tables;
DROP TABLE IF EXISTS dish_snapshots;
DROP TABLE IF EXISTS dishes;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS question_models;
DROP TABLE IF EXISTS paragraph_content_models;
DROP TABLE IF EXISTS passage_models;
DROP TABLE IF EXISTS section_models;
DROP TABLE IF EXISTS reading_test_models;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS question_type;

-- Recreate types
CREATE TYPE question_type AS ENUM ('MultipleChoice', 'TrueFalseNotGiven', 'Matching', 'ShortAnswer');

-- Recreate tables
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'Employee',
    phone VARCHAR(20),
    image VARCHAR(255),
    address TEXT,
    favorite_food INTEGER[] DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_role ON users(role);

ALTER TABLE users
ADD CONSTRAINT check_valid_role
CHECK (role IN ('Admin', 'Employee', 'Manager', 'Guest'));

CREATE TABLE sessions (
    id varchar(255) PRIMARY KEY NOT NULL,
    user_email VARCHAR(100) NOT NULL,
    refresh_token TEXT NOT NULL,
    is_revoked BOOLEAN DEFAULT false,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE reading_test_models (
    id SERIAL PRIMARY KEY,
    test_number INTEGER NOT NULL,
    sections JSONB NOT NULL
);

CREATE TABLE section_models (
    id SERIAL PRIMARY KEY,
    section_number INTEGER NOT NULL,
    time_allowed INTEGER NOT NULL,
    passages JSONB NOT NULL
);

CREATE TABLE passage_models (
    id SERIAL PRIMARY KEY,
    passage_number INTEGER NOT NULL,
    title TEXT NOT NULL,
    content JSONB NOT NULL,
    questions JSONB NOT NULL
);

CREATE TABLE paragraph_content_models (
    id SERIAL PRIMARY KEY,
    paragraph_summary TEXT NOT NULL,
    key_words TEXT NOT NULL,
    key_sentence TEXT NOT NULL
);

CREATE TABLE question_models (
    id SERIAL PRIMARY KEY,
    question_number INTEGER NOT NULL,
    type question_type NOT NULL,
    content TEXT NOT NULL,
    options TEXT[] NULL,
    correct_answer JSONB NULL
);

CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    author_id BIGINT NOT NULL,
    parent_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    avatar VARCHAR(255),
    role VARCHAR(50) DEFAULT 'Employee',
    owner_id BIGINT,
    favorite_sets BIGINT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE dishes (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    description TEXT NOT NULL,
    image VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'Available',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE dish_snapshots (
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

CREATE TABLE tables (
    number INTEGER PRIMARY KEY,
    capacity INTEGER NOT NULL,
    status VARCHAR(50) DEFAULT 'Available',
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE guests (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    table_number INTEGER,
    refresh_token VARCHAR(255),
    refresh_token_expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (table_number) REFERENCES tables(number) ON DELETE SET NULL
);

CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    guest_id BIGINT,
    table_number INTEGER,
    dish_snapshot_id BIGINT NOT NULL,
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

CREATE TABLE refresh_tokens (
    token VARCHAR(255) PRIMARY KEY,
    account_id BIGINT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE TABLE sockets (
    socket_id VARCHAR(255) PRIMARY KEY,
    account_id BIGINT UNIQUE,
    guest_id BIGINT UNIQUE,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE SET NULL,
    FOREIGN KEY (guest_id) REFERENCES guests(id) ON DELETE SET NULL
);

-- Updated table for sets
CREATE TABLE sets (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    user_id BIGINT,
    is_favourite BOOLEAN DEFAULT FALSE,
    like_by BIGINT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_public BOOLEAN DEFAULT FALSE,
    image VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- New table for set_dishes (junction table between sets and dishes)
CREATE TABLE set_dishes (
    id BIGSERIAL PRIMARY KEY,
    set_id BIGINT NOT NULL,
    dish_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (set_id) REFERENCES sets(id) ON DELETE CASCADE,
    FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE CASCADE,
    UNIQUE (set_id, dish_id)
);

-- Updated table for set snapshots
CREATE TABLE set_snapshots (
    id BIGSERIAL PRIMARY KEY,
    original_set_id BIGINT,
    set_id BIGINT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    user_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_public BOOLEAN DEFAULT FALSE,
    image VARCHAR(255),
    FOREIGN KEY (original_set_id) REFERENCES sets(id) ON DELETE SET NULL,
    FOREIGN KEY (set_id) REFERENCES sets(id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- New table for set_snapshot_dishes
CREATE TABLE set_snapshot_dishes (
    id BIGSERIAL PRIMARY KEY,
    set_snapshot_id BIGINT NOT NULL,
    dish_snapshot_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (set_snapshot_id) REFERENCES set_snapshots(id) ON DELETE CASCADE,
    FOREIGN KEY (dish_snapshot_id) REFERENCES dish_snapshots(id) ON DELETE CASCADE,
    UNIQUE (set_snapshot_id, dish_snapshot_id)
);

-- Add constraints and indexes
ALTER TABLE accounts
ADD CONSTRAINT fk_owner
FOREIGN KEY (owner_id) REFERENCES accounts(id) ON DELETE SET NULL;

CREATE INDEX idx_sets_user_id ON sets(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_sets_is_favourite ON sets(is_favourite);
CREATE INDEX idx_sets_like_by ON sets USING GIN (like_by);
CREATE INDEX idx_set_snapshots_original_set_id ON set_snapshots(original_set_id);

-- New indexes for is_public column
CREATE INDEX idx_sets_is_public ON sets(is_public);
CREATE INDEX idx_set_snapshots_is_public ON set_snapshots(is_public);



-- set 
