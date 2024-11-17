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

	GetNodeQuery = `
		WITH relation_data AS (
			SELECT
				nd.id AS node_id,
				COALESCE(
								json_agg(DISTINCT rl.relative_id) 
									FILTER (WHERE rl.relation_type = 'Родитель' AND NOT rl.is_deleted),
								'[]'::json
				) AS parent_ids,
				COALESCE(
								json_agg(DISTINCT rl.relative_id) 
									FILTER (WHERE rl.relation_type = 'Супруг' AND NOT rl.is_deleted),
								'[]'::json
				) AS spouse_ids
			FROM public.node nd
					 LEFT JOIN public.relation rl ON rl.node_id = nd.id
			WHERE nd.id = $1
			GROUP BY nd.id
		),
			 children_data AS(
				 SELECT
					 nd.id AS node_id,
					 COALESCE(
									 json_agg(DISTINCT rl.node_id) 
										FILTER (WHERE rl.relation_type = 'Родитель' AND NOT rl.is_deleted),
									 '[]'::json
					 ) AS children_ids
				 FROM public.node nd
						  LEFT JOIN public.relation rl ON rl.relative_id = nd.id
				 WHERE nd.id = $1
				 GROUP BY nd.id
			 ),
			 tree_data AS (
				 SELECT
					 nd.id AS node_id,
					 nd.layer_id,
					 nd.name,
					 nd.birthdate,
					 nd.deathdate,
					 nd.preview_path,
					 nd.is_spouse,
					 nd.gender
				 FROM public.tree tr
						  JOIN public.layer lr ON lr.tree_id = tr.id
						  JOIN public.node nd ON lr.id = nd.layer_id
				 WHERE NOT nd.is_deleted AND nd.id = $1
			 )
		SELECT
			td.*,
			COALESCE(jsonb_build_object(
							 'children', children_data.children_ids,
							 'parents', relation_data.parent_ids,
							 'spouses', relation_data.spouse_ids
					 ), '{}'::jsonb) AS relation_json
		FROM tree_data td
				 LEFT JOIN relation_data ON td.node_id = relation_data.node_id
				 LEFT JOIN children_data ON td.node_id = children_data.node_id;
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

	UpdateGenderQuery = `
		UPDATE public.node
		SET gender = $1
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

	DeleteRelationsQuery = `
		UPDATE public.relation
		SET is_deleted = TRUE
		WHERE relative_id = $1 OR node_id = $1;
		`
)
