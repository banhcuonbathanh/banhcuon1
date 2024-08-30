-- Step 1: Create a new UUID column
ALTER TABLE reading_res_models ADD COLUMN new_id UUID DEFAULT gen_random_uuid();

-- Step 2: Copy data from the old id to the new_id (if needed)
UPDATE reading_res_models SET new_id = id::text::uuid;

-- Step 3: Drop the primary key constraint
ALTER TABLE reading_res_models DROP CONSTRAINT reading_res_models_pkey;

-- Step 4: Drop the old id column
ALTER TABLE reading_res_models DROP COLUMN id;

-- Step 5: Rename the new_id column to id
ALTER TABLE reading_res_models RENAME COLUMN new_id TO id;

-- Step 6: Set the new id column as primary key
ALTER TABLE reading_res_models ADD PRIMARY KEY (id);

-- Step 7: Set default for new rows
ALTER TABLE reading_res_models ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- Step 8: Drop the sequence if it exists
DROP SEQUENCE IF EXISTS reading_res_models_id_seq;