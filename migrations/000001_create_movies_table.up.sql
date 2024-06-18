CREATE TABLE IF NOT EXISTS movies (
    id bigserial PRIMARY KEY, 
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), 
    title text NOT NULL, 
    year integer NOT NULL, 
    runtime integer NOT NULL, 
    genres text[] NOT NULL, 
    version integer NOT NULL DEFAULT 1
);

INSERT INTO movies (title, year, runtime, genres) 
VALUES 
    ('Moana', 2018, 134, ARRAY['action', 'adventure']), 
    ('Deadpool', 2016, 108, ARRAY['action', 'comedy', 'superhero']), 
    ('The Breakfast Club', 1986, 96, ARRAY['drama']);