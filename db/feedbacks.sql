CREATE TABLE feedbacks(
	id bigserial PRIMARY KEY,
	message text NOT NULL,
	useruuid uuid NOT NULL,
	created timestamp without time zone NOT NULL
);