CREATE TYPE "story_status" AS ENUM (
  'draft',
  'published',
  'archived'
);

CREATE TABLE "users" (
  "id" varchar PRIMARY KEY,
  "name" varchar NOT NULL,
  "bio" text NOT NULL,
  "profile_image" varchar NOT NULL DEFAULT '',
  "location" varchar NOT NULL,
  "website" varchar NOT NULL DEFAULT '',
  "followers" integer NOT NULL DEFAULT 0,
  "following" integer NOT NULL DEFAULT 0,
  "stories_count" integer NOT NULL DEFAULT 0,
  "is_verified" boolean NOT NULL DEFAULT false
);


CREATE TABLE "stories" (
  "id" varchar PRIMARY KEY,
  "title" varchar NOT NULL,
  "description" text,
  "cover_image" varchar,
  "author_id" varchar NOT NULL,
  "likes" integer NOT NULL DEFAULT 0,
  "views" integer NOT NULL DEFAULT 0,
  "published_date" timestamptz,
  "last_edited" timestamptz NOT NULL,
  "story_type" varchar NOT NULL,
  "status" story_status NOT NULL,
  "genres" text[]
);

CREATE TABLE "chapters" (
  "id" varchar PRIMARY KEY,
  "story_id" varchar NOT NULL,
  "title" varchar NOT NULL,
  "content" text NOT NULL,
  "is_complete" boolean NOT NULL DEFAULT false
);

CREATE TABLE "story_collaborators" (
  "story_id" varchar NOT NULL,
  "user_id" varchar NOT NULL,
  PRIMARY KEY ("story_id", "user_id")
);

ALTER TABLE "stories" ADD FOREIGN KEY ("author_id") REFERENCES "users" ("id");

ALTER TABLE "chapters" ADD FOREIGN KEY ("story_id") REFERENCES "stories" ("id") ON DELETE CASCADE;

ALTER TABLE "story_collaborators" ADD FOREIGN KEY ("story_id") REFERENCES "stories" ("id") ON DELETE CASCADE;

ALTER TABLE "story_collaborators" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
-- Step 1: Drop the existing foreign key constraint
ALTER TABLE story_collaborators
DROP CONSTRAINT IF EXISTS story_collaborators_user_id_fkey;

-- Step 2: Re-add the foreign key with ON DELETE CASCADE
ALTER TABLE story_collaborators
ADD CONSTRAINT story_collaborators_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
