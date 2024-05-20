CREATE TABLE users (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    fullname VARCHAR(255),
    password VARCHAR(255)
);

CREATE TABLE classes (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(255),
    lecturer VARCHAR(255),
    description TEXT,
    icon VARCHAR(255),
    code VARCHAR(6)
);

CREATE TABLE tasks (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    title VARCHAR(255),
    description TEXT,
    detail TEXT,
    submission VARCHAR(255),
    task_type VARCHAR(255),
    deadline TIMESTAMP
);

CREATE TABLE user_classes (
    user_id UUID REFERENCES users(id),
    class_id UUID REFERENCES classes(id),
    PRIMARY KEY (user_id, class_id)
);

CREATE TABLE class_tasks (
    class_id UUID REFERENCES classes(id),
    task_id UUID REFERENCES tasks(id),
    PRIMARY KEY (class_id, task_id)
);

CREATE TABLE user_tasks (
    user_id UUID REFERENCES users(id),
    task_id UUID REFERENCES tasks(id),
    status BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (user_id, task_id)
);
