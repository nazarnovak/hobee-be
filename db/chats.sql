CREATE TABLE chats(
	id bigserial PRIMARY KEY,
	user1 uuid NOT NULL,
	user2 uuid NOT NULL,
	room uuid NOT NULL,
	messages bytea NOT NULL,
	started timestamp without time zone NOT NULL,
	finished timestamp without time zone NOT NULL,
	user1_likes text[],
	user1_reports text[],
	user2_likes text[],
	user2_reports text[]
);