CREATE TABLE user (  
    user_id BINARY(32) NOT NULL,
    create_time DATETIME,
    username VARCHAR(255),
    password VARCHAR(255),
    PRIMARY KEY (user_id)
);