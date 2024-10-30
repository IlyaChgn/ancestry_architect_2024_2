package repository

const (
	CheckPermissionForNodeQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM public.node nd
				 JOIN public.layer lr ON nd.layer_id = lr.id
				 JOIN public.tree tr ON tr.id = lr.tree_id
				 JOIN public.permission pm ON pm.tree_id = tr.id
			WHERE pm.user_id = $2 AND nd.id = $1
		
			UNION ALL
		
			SELECT 1
			FROM public.node nd
				 JOIN public.layer lr ON nd.layer_id = lr.id
				 JOIN public.tree tr ON tr.id = lr.tree_id
			WHERE tr.user_id = $2
		);
		`

	GetRelativeNodeQuery = `
		SELECT nd.id, lr.number, lr.id
		FROM public.node nd
			JOIN public.layer lr ON nd.layer_id = lr.id
		WHERE nd.id = $1;
		`

	GetParentsQuery = `
		SELECT relative_id
		FROM public.relation
		WHERE node_id = $1 AND relation_type = 'Родитель';
		`

	GetDescriptionQuery = `
		SELECT id, description
		FROM public.description
		WHERE node_id = $1;
		`

	GetPreviewQuery = `
		SELECT preview_path
		FROM public.node
		WHERE id = $1;
		`

	GetLayerQuery = `
		SELECT id
		FROM public.layer
		WHERE tree_id = $1 AND number = $2;
		`

	CreateNodeQuery = `
		INSERT 
		INTO public.node (layer_id, name, is_spouse, gender) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, layer_id, name, birthdate, deathdate, is_spouse, gender;
		`

	UpdateBirthdateQuery = `
		UPDATE public.node
		SET birthdate = $1
		WHERE id = $2
		RETURNING birthdate;
		`

	UpdateDeathdateQuery = `
		UPDATE public.node
		SET deathdate = $1
		WHERE id = $2
		RETURNING deathdate;
		`

	UpdatePreviewQuery = `
		UPDATE public.node
		SET preview_path = $1
		WHERE id = $2
		RETURNING preview_path;
		`

	InsertDescriptionQuery = `
		INSERT
		INTO public.description (description, node_id)
		VALUES ($1, $2);
		`

	UpdateDescriptionQuery = `
		UPDATE public.description
		SET description = $1
		WHERE node_id = $2;
		`

	UpdateNameQuery = `
		UPDATE public.node
		SET name = $1
		WHERE id = $2;
		`

	SetRelativeQuery = `
		INSERT
		INTO public.relation (relative_id, node_id, relation_type)
		VALUES ($1, $2, $3);
		`

	CreateLayerQuery = `
		INSERT
		INTO public.layer (tree_id, number) 
		VALUES ($1, $2)
		RETURNING id; 
		`

	DeleteNodeQuery = `
		UPDATE public.node
		SET is_deleted = TRUE
		WHERE id = $1;
		`
)
