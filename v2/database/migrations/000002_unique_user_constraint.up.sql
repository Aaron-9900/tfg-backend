ALTER TABLE users
ADD CONSTRAINT user_unique UNIQUE(name, email);
