-- Add new columns
ALTER TABLE orders
ADD COLUMN bow_chili BIGINT DEFAULT 0,
ADD COLUMN bow_no_chili BIGINT DEFAULT 0;

-- Modify existing columns
ALTER TABLE orders
ALTER COLUMN guest_id DROP NOT NULL,
ALTER COLUMN user_id DROP NOT NULL,
ALTER COLUMN table_number TYPE BIGINT,
ALTER COLUMN total_price TYPE INTEGER;

-- Update constraints
ALTER TABLE orders
DROP CONSTRAINT IF EXISTS guest_or_user_check,
ADD CONSTRAINT guest_or_user_check CHECK (
    (is_guest = TRUE AND guest_id IS NOT NULL AND user_id IS NULL) OR
    (is_guest = FALSE AND user_id IS NOT NULL AND guest_id IS NULL)
);

-- Update foreign key for table_number
ALTER TABLE orders
DROP CONSTRAINT IF EXISTS orders_table_number_fkey,
ADD CONSTRAINT orders_table_number_fkey
    FOREIGN KEY (table_number)
    REFERENCES tables(number)
    ON DELETE SET NULL;

-- Update indexes
DROP INDEX IF EXISTS idx_orders_table_number;
CREATE INDEX idx_orders_table_number ON orders(table_number) WHERE table_number IS NOT NULL;