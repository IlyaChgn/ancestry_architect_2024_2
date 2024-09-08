package repository

const (
	GetProfileByIDQuery = `
		SELECT id, user_id, name, surname, birthdate, gender, avatar_path
		FROM public.profile
		WHERE user_id = $1;
		`

	CreateProfileQuery = `
		INSERT 
		INTO public.profile (user_id)
		VALUES ($1)
		RETURNING id, user_id;
		`
)
