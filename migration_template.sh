mkdir -p ./migrations
cat << EOF > ./migrations/$1.up.sql
CREATE TABLE IF NOT EXISTS ${2%%.*} (
	uuid uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	/*----------------------------*/
    -- extras

);
EOF

cat << EOF > ./migrations/$1.down.sql
DROP TABLE IF EXISTS ${2%%.*};
EOF