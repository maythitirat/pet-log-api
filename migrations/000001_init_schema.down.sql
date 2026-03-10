-- 000001_init_schema.down.sql
-- Rollback initial database schema

DROP TRIGGER IF EXISTS update_pets_updated_at ON pets;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS pets;
DROP TABLE IF EXISTS users;
