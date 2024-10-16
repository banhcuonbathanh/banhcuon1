-- Modify the sets table
ALTER TABLE sets
ADD COLUMN is_public BOOLEAN DEFAULT FALSE;

-- Modify the set_snapshots table
ALTER TABLE set_snapshots
ADD COLUMN is_public BOOLEAN DEFAULT FALSE;

-- Create an index for the new is_public column
CREATE INDEX idx_sets_is_public ON sets(is_public);
CREATE INDEX idx_set_snapshots_is_public ON set_snapshots(is_public);

-- Remove existing constraints
ALTER TABLE sets
DROP CONSTRAINT IF EXISTS check_user_id_positive;

ALTER TABLE set_snapshots
DROP CONSTRAINT IF EXISTS check_set_snapshots_user_id_positive;