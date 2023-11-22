CREATE TABLE "sessions" (
    "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" BOOLEAN NOT NULL DEFAULT FALSE,
  "expires_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);


ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");