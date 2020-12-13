CREATE TABLE chats(
	id bigserial PRIMARY KEY,
	user1 uuid NOT NULL,
	user2 uuid NOT NULL,
	room uuid NOT NULL,
	messages bytea NOT NULL,
	started timestamp without time zone NOT NULL,
	finished timestamp without time zone NOT NULL,
	user1_like character varying(16),
	user1_report character varying(16),
	user2_like character varying(16),
	user2_report character varying(16)
);