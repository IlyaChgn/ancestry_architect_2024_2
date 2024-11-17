CREATE TYPE relations AS ENUM ('Супруг', 'Родитель');

ALTER TABLE public.relation
RENAME COLUMN parent_id TO relative_id;

ALTER TABLE public.relation
RENAME COLUMN child_id TO node_id;

ALTER TABLE public.relation
ADD COLUMN IF NOT EXISTS relation_type relations
    NOT NULL
    DEFAULT 'Родитель';