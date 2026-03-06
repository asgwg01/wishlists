DROP TRIGGER auto_update_users_trig ON users;
DROP INDEX email_users_idx;
DROP TABLE users;
DROP FUNCTION auto_update_time();