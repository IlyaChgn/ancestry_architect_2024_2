CREATE TABLE IF NOT EXISTS public.tree
(
    id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    -- Owner`s ID
    user_id BIGINT NOT NULL
        REFERENCES public.user (id),
    -- Tree`s name
    name TEXT NOT NULL
        CHECK (name <> '')
        CONSTRAINT max_len_name CHECK (LENGTH(name) <= 100)
);

CREATE TABLE IF NOT EXISTS public.layer
(
    id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    tree_id BIGINT NOT NULL
        REFERENCES public.tree (id),
    -- Layer`s number
    number SMALLINT NOT NULL
        CONSTRAINT positive_num CHECK (number >= 0)
);

CREATE TABLE IF NOT EXISTS public.node
(
    id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    layer_id BIGINT NOT NULL
        REFERENCES public.layer (id),
    -- Ancestor`s full name
    name TEXT NOT NULL
        CHECK (name <> '')
        CONSTRAINT max_len_name CHECK (LENGTH(name) <= 100),
    birthdate DATE DEFAULT NULL,
    deathdate DATE DEFAULT NULL,
    -- Preview image
    preview_path TEXT DEFAULT NULL
        CONSTRAINT max_len_preview_path CHECK (LENGTH(preview_path) <= 256)
);

CREATE TABLE IF NOT EXISTS public.relation
(
    id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    parent_id BIGINT NOT NULL
        REFERENCES public.node (id),
    child_id BIGINT NOT NULL
        REFERENCES public.node (id)
);

CREATE TABLE IF NOT EXISTS description
(
    id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    node_id BIGINT UNIQUE NOT NULL
        REFERENCES public.node (id),
    description TEXT NOT NULL
        CHECK (description <> '')
        CONSTRAINT max_len_description CHECK (LENGTH(description) <= 500)
)
