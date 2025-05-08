CREATE TABLE "users"
(
    "id"                BIGSERIAL    NOT NULL,
    "name"              VARCHAR(255) NOT NULL,
    "email"             VARCHAR(255) NOT NULL,
    "email_verified_at" TIMESTAMP(0),
    "password"          VARCHAR(255) NOT NULL,
    "created_at"        TIMESTAMP(0) NOT NULL,
    "updated_at"        TIMESTAMP(0) NOT NULL,
    CONSTRAINT "pidx__users__id" PRIMARY KEY ("id"),
    CONSTRAINT "uidx__users__email" UNIQUE ("email")
);