ALTER TABLE
    tasks
ADD
    COLUMN class_id uuid,
ADD
    CONSTRAINT fk_classes_tasks FOREIGN KEY (class_id) REFERENCES classes(id);