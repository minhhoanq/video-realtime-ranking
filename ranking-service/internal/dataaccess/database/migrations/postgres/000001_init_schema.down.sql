DROP TABLE IF EXISTS "video";
DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS "category";

-- Xóa các ENUM type
DROP TYPE IF EXISTS event_type;
DROP TYPE IF EXISTS user_role;

-- (Tùy chọn) Xóa extension uuid-ossp
DROP EXTENSION IF EXISTS "uuid-ossp";