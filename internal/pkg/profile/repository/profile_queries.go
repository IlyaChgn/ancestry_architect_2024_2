package repository

const (
	GetProfileByIDQuery = `
		SELECT id, user_id, name, surname, birthdate, gender, avatar_path
		FROM public.profile
		WHERE user_id = $1;
		`

	GetAvatarQuery = `
		SELECT avatar_path
		FROM public.profile
		WHERE user_id = $1;
		`

	CreateProfileQuery = `
		INSERT 
		INTO public.profile (user_id)
		VALUES ($1)
		RETURNING id, user_id;
		`

	UpdateAvatarQuery = `
		UPDATE public.profile
		SET avatar_path = $1
		WHERE user_id = $2;
		`

	UpdateBirthdateQuery = `
		UPDATE public.profile
		SET birthdate = $1
		WHERE user_id = $2;
		`

	UpdateGenderQuery = `
		UPDATE public.profile
		SET gender = $1
		WHERE user_id = $2;
		`

	UpdateProfileQuery = `
		UPDATE public.profile
		SET name = $1,
			surname = $2
		WHERE user_id = $3
		RETURNING id, user_id, name, surname, birthdate, gender, avatar_path; 
		`
)
