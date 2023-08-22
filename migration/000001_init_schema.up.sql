CREATE TABLE "users" (
  "id" BIGSERIAL PRIMARY KEY,
  "username" VARCHAR(255) UNIQUE NOT NULL,
  "full_name" VARCHAR(255) NOT NULL,
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "salt" VARCHAR(255) NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01',
  "creator_username" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "sessions" (
  "id" bigserial PRIMARY KEY,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" VARCHAR(255) NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "owner_id" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "videos" (
  "id" BIGSERIAL PRIMARY KEY,
  "title" VARCHAR(255) NOT NULL,
  "description" TEXT,
  "thumnail_url" varchar,
  "status" VARCHAR(24) NOT NULL,
  "old_status" VARCHAR(24),
  "owner_id" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "video_files" (
  "id" BIGSERIAL PRIMARY KEY,
  "video_id" bigserial NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  "path" VARCHAR(255) NOT NULL,
  "format" VARCHAR(10) NOT NULL,
  "provider" VARCHAR(10) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "categories" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR(255) UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "video_categories" (
  "id" BIGSERIAL PRIMARY KEY,
  "video_id" bigserial NOT NULL,
  "category_id" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "comments" (
  "id" BIGSERIAL PRIMARY KEY,
  "video_id" bigserial NOT NULL,
  "user_id" bigserial NOT NULL,
  "content" TEXT NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE INDEX ON "sessions" ("owner_id");

CREATE INDEX ON "videos" ("owner_id");

CREATE INDEX ON "video_files" ("video_id");

CREATE INDEX ON "video_categories" ("video_id");

CREATE INDEX ON "video_categories" ("category_id");

CREATE UNIQUE INDEX ON "video_categories" ("video_id", "category_id");

CREATE INDEX ON "comments" ("video_id");

CREATE INDEX ON "comments" ("user_id");

ALTER TABLE "users" ADD FOREIGN KEY ("creator_username") REFERENCES "users" ("username");

ALTER TABLE "sessions" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id");

ALTER TABLE "videos" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id");

ALTER TABLE "video_files" ADD FOREIGN KEY ("video_id") REFERENCES "videos" ("id");

ALTER TABLE "video_categories" ADD FOREIGN KEY ("video_id") REFERENCES "videos" ("id");

ALTER TABLE "video_categories" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("video_id") REFERENCES "videos" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
