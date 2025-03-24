-- Create stories table
CREATE TABLE stories (
    story_id SERIAL PRIMARY KEY,
    originalstory TEXT NOT NULL,
    pulledrequests INT DEFAULT 0,
    updatedstory TEXT NOT NULL,
    author_id INT REFERENCES authors(id) ON DELETE CASCADE
);
