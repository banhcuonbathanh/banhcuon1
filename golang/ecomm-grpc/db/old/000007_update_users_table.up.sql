-- Modify reading_res_models table
ALTER TABLE reading_res_models
    ALTER COLUMN id TYPE UUID USING (gen_random_uuid()),
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- If you want to drop the sequence created by BIGSERIAL (optional)
DROP SEQUENCE IF EXISTS reading_res_models_id_seq;