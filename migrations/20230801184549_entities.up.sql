CREATE TABLE IF NOT EXISTS entities (
	uuid uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	/*----------------------------*/
    -- extras
    name varchar,
    email_address varchar,
    username varchar,
    password varchar,
    icon varchar,
    description varchar,
    type int,
    user_uuid uuid REFERENCES users(uuid)
);
