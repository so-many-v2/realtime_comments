CREATE TABLE appuser (
     id       BIGSERIAL    PRIMARY KEY,
     username VARCHAR(128) NOT NULL
);

CREATE TABLE posts (
    id        BIGSERIAL    PRIMARY KEY,
    user_id   BIGINT       NOT NULL REFERENCES appuser(id) ON DELETE CASCADE,
    title     VARCHAR(256) NOT NULL,
    text      TEXT         NOT NULL,
    created   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated   TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE comment (
    id        BIGSERIAL   PRIMARY KEY,
    user_id   BIGINT      NOT NULL REFERENCES appuser(id) ON DELETE CASCADE,
    text      TEXT        NOT NULL,
    created   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE post_comments (
    id         BIGSERIAL PRIMARY KEY,
    post_id    BIGINT NOT NULL REFERENCES posts(id)   ON DELETE CASCADE,
    comment_id BIGINT NOT NULL REFERENCES comment(id) ON DELETE CASCADE
);