CREATE TABLE IF NOT EXISTS users (
     uuid uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    /*----------------------------*/
    -- extras
     email_address varchar,
     password varchar,
     active boolean
);

CREATE UNIQUE INDEX users_unique_email_address on users(email_address) where active is true;
