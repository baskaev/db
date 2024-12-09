-- CREATE TABLE movies (
--     id SERIAL PRIMARY KEY,
--     title VARCHAR(255) NOT NULL,
--     description TEXT,
--     release_date DATE
-- );

CREATE TABLE movies (
    code VARCHAR(255) PRIMARY KEY, 
    title VARCHAR(255) NOT NULL,   
    rating VARCHAR(10),            
    year VARCHAR(4),               
    image_link TEXT                
);

