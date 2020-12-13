CREATE TABLE contacts(
	id bigserial PRIMARY KEY,
	name text NOT NULL,
	email character varying(512) NOT NULL,
	message text NOT NULL,
	useruuid uuid NOT NULL
);