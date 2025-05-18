BEGIN TRANSACTION;

CREATE TABLE "words"
(
    "id"         SERIAL       NOT NULL,
    word         VARCHAR(5)   NOT NULL,
    "created_at" TIMESTAMP(0) NOT NULL,
    CONSTRAINT "pidx__words__id" PRIMARY KEY ("id"),
    CONSTRAINT "uidx__words__word" UNIQUE ("word")
);

CREATE TABLE "games"
(
    "id"          BIGSERIAL    NOT NULL,
    "user_id"     BIGINT       NOT NULL,
    "word_id"     INT          NOT NULL,
    "guess_limit" SMALLINT     NOT NULL,
    "is_playing"  BOOLEAN      NOT NULL DEFAULT TRUE,
    "is_won"      BOOLEAN,
    "created_at"  TIMESTAMP(0) NOT NULL,
    "updated_at"  TIMESTAMP(0) NOT NULL,
    CONSTRAINT "pidx__games__id" PRIMARY KEY ("id"),
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT
        NOT DEFERRABLE INITIALLY IMMEDIATE,
    FOREIGN KEY ("word_id") REFERENCES "words" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT
        NOT DEFERRABLE INITIALLY IMMEDIATE
);
CREATE INDEX "idx__games__user_id__is_playing" ON "games" ("user_id", "is_playing");

CREATE TABLE "guesses"
(
    "id"         BIGSERIAL    NOT NULL,
    "game_id"    BIGINT       NOT NULL,
    "word_id"    INT          NOT NULL,
    "created_at" TIMESTAMP(0) NOT NULL,
    CONSTRAINT "pidx__guesses__id" PRIMARY KEY ("id"),
    FOREIGN KEY ("game_id") REFERENCES "games" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT
        NOT DEFERRABLE INITIALLY IMMEDIATE,
    FOREIGN KEY ("word_id") REFERENCES "words" ("id") ON DELETE RESTRICT ON UPDATE RESTRICT
        NOT DEFERRABLE INITIALLY IMMEDIATE
);

COMMIT;
