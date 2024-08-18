CREATE OR REPLACE FUNCTION change_profile_update_time()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_time := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_profile_trigger
BEFORE UPDATE ON public.profile
FOR EACH ROW
EXECUTE PROCEDURE change_profile_update_time();