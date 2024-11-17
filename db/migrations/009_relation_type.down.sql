ALTER TABLE public.relation
RENAME COLUMN relative_id TO parent_id;

ALTER TABLE public.relation
RENAME COLUMN node_id TO child_id;

ALTER TABLE public.relation
DROP COLUMN IF EXISTS relation_type;

DROP TYPE IF EXISTS relations;