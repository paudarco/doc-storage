BEGIN;

DROP INDEX IF EXISTS idx_documents_user_id;
DROP INDEX IF EXISTS idx_documents_name;

DROP INDEX IF EXISTS idx_users_login;
DROP INDEX IF EXISTS idx_users_id;

DROP TABLE IF EXISTS documents;

DROP TABLE IF EXISTS users;

COMMIT;