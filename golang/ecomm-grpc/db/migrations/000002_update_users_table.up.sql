-- Create ENUM type
CREATE TYPE question_type AS ENUM ('MultipleChoice', 'TrueFalseNotGiven', 'Matching', 'ShortAnswer');

-- Create tables
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
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_valid_role CHECK (role IN ('Admin', 'Employee', 'Manager', 'Guest'))
);

CREATE TABLE tables (
    number BIGINT PRIMARY KEY,
    capacity INTEGER NOT NULL,
    status VARCHAR(50) DEFAULT 'Available',
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE guests (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    table_number BIGINT,
    refresh_token VARCHAR(255),
    refresh_token_expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (table_number) REFERENCES tables(number) ON DELETE SET NULL
);

CREATE TABLE guest_sessions (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    refresh_token TEXT NOT NULL,
    is_revoked BOOLEAN NOT NULL DEFAULT false,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
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
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner_id) REFERENCES accounts(id) ON DELETE SET NULL
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

CREATE TABLE sets (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    user_id BIGINT,
    price INTEGER NOT NULL DEFAULT 0,
    is_favourite BOOLEAN DEFAULT FALSE,
    like_by BIGINT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_public BOOLEAN DEFAULT FALSE,
    image VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE set_dishes (
    id BIGSERIAL PRIMARY KEY,
    set_id BIGINT NOT NULL,
    dish_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (set_id) REFERENCES sets(id) ON DELETE CASCADE,
    FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE CASCADE,
    UNIQUE (set_id, dish_id)
);

CREATE TABLE set_snapshots (
    id BIGSERIAL PRIMARY KEY,
    original_set_id BIGINT,
    set_id BIGINT,
    price INTEGER NOT NULL DEFAULT 0,
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

CREATE TABLE set_snapshot_dishes (
    id BIGSERIAL PRIMARY KEY,
    set_snapshot_id BIGINT NOT NULL,
    dish_snapshot_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (set_snapshot_id) REFERENCES set_snapshots(id) ON DELETE CASCADE,
    FOREIGN KEY (dish_snapshot_id) REFERENCES dish_snapshots(id) ON DELETE CASCADE,
    UNIQUE (set_snapshot_id, dish_snapshot_id)
);

CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    guest_id BIGINT,
    user_id BIGINT,
    is_guest BOOLEAN NOT NULL,
    table_number BIGINT,
    order_handler_id BIGINT,
    status VARCHAR(50) DEFAULT 'Pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    total_price INTEGER,
    topping VARCHAR(255),
    tracking_order VARCHAR(255),
    take_away BOOLEAN NOT NULL DEFAULT false,
    chili_number BIGINT DEFAULT 0,
    table_token VARCHAR(255) NOT NULL,
    order_name VARCHAR(255),
    FOREIGN KEY (guest_id) REFERENCES guests(id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (table_number) REFERENCES tables(number) ON DELETE SET NULL,
    FOREIGN KEY (order_handler_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT guest_or_user_check CHECK (
        (is_guest = TRUE AND guest_id IS NOT NULL AND user_id IS NULL) OR
        (is_guest = FALSE AND user_id IS NOT NULL AND guest_id IS NULL)
    )
);

CREATE TABLE dish_order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    dish_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE CASCADE
);

CREATE TABLE set_order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    set_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (set_id) REFERENCES sets(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_sets_user_id ON sets(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_sets_is_favourite ON sets(is_favourite);
CREATE INDEX idx_sets_like_by ON sets USING GIN (like_by);
CREATE INDEX idx_sets_is_public ON sets(is_public);
CREATE INDEX idx_orders_guest_id ON orders(guest_id) WHERE guest_id IS NOT NULL;
CREATE INDEX idx_orders_user_id ON orders(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_orders_order_handler_id ON orders(order_handler_id) WHERE order_handler_id IS NOT NULL;
CREATE INDEX idx_orders_table_number ON orders(table_number) WHERE table_number IS NOT NULL;
CREATE INDEX idx_orders_is_guest ON orders(is_guest);
CREATE INDEX idx_guest_sessions_refresh_token ON guest_sessions(refresh_token);




-- Create tables delivery order

CREATE TABLE deliveries (
    id BIGSERIAL PRIMARY KEY,
    guest_id BIGINT,
    user_id BIGINT,
    is_guest BOOLEAN NOT NULL,
    table_number BIGINT,
    order_handler_id BIGINT,
    status VARCHAR(50) DEFAULT 'Pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    total_price INTEGER,
    topping VARCHAR(255),
    bow_no_chili BIGINT DEFAULT 0,
    take_away BOOLEAN NOT NULL DEFAULT false,
    chili_number BIGINT DEFAULT 0,
    table_token VARCHAR(255) NOT NULL,
    client_name VARCHAR(255),
    -- Suggested additional fields
    delivery_address TEXT,
    delivery_contact VARCHAR(100),
    delivery_notes TEXT,
    scheduled_time TIMESTAMP WITH TIME ZONE,
    delivery_fee INTEGER DEFAULT 0,
    driver_id BIGINT,
    delivery_status VARCHAR(50) DEFAULT 'Pending',
    estimated_delivery_time TIMESTAMP WITH TIME ZONE,
    actual_delivery_time TIMESTAMP WITH TIME ZONE,
    
    FOREIGN KEY (guest_id) REFERENCES guests(id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (table_number) REFERENCES tables(number) ON DELETE SET NULL,
    FOREIGN KEY (order_handler_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (driver_id) REFERENCES users(id) ON DELETE SET NULL,
    
    CONSTRAINT guest_or_user_check CHECK (
        (is_guest = TRUE AND guest_id IS NOT NULL AND user_id IS NULL) OR
        (is_guest = FALSE AND user_id IS NOT NULL AND guest_id IS NULL)
    ),
    CONSTRAINT valid_delivery_status CHECK (
        delivery_status IN ('Pending', 'Assigned', 'Picked Up', 'In Transit', 'Delivered', 'Failed', 'Cancelled')
    )
);

CREATE TABLE dish_delivery_items (
    id BIGSERIAL PRIMARY KEY,
    delivery_id BIGINT NOT NULL,
    dish_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL,
    FOREIGN KEY (delivery_id) REFERENCES deliveries(id) ON DELETE CASCADE,
    FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE CASCADE
);

-- Index for faster lookups
CREATE INDEX idx_deliveries_client_name ON deliveries(client_name);
CREATE INDEX idx_deliveries_status ON deliveries(status);
CREATE INDEX idx_deliveries_delivery_status ON deliveries(delivery_status);