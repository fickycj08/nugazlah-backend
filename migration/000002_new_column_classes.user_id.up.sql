ALTER TABLE
    classes
ADD
    COLUMN user_id uuid,
ADD
    CONSTRAINT fk_classes_users FOREIGN KEY (user_id) REFERENCES users(id);