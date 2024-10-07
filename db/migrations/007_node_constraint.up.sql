ALTER TABLE public.node
ADD CONSTRAINT deathdate_after_birthdate CHECK(deathdate > birthdate);