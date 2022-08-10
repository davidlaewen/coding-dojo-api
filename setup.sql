-- Set up `tasks` table

CREATE TABLE IF NOT EXISTS tasks (
    uuid text PRIMARY KEY,
    title text NOT NULL,
    description text
);

-- Create index for efficient search by UUID
CREATE INDEX IF NOT EXISTS tasks_uuid ON tasks(uuid);
