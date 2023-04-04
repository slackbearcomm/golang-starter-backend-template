package resolvers

/////////////////
//   Queries   //
/////////////////

// func (r *queryResolver) GetObject(ctx context.Context, uid string) (*graph.GetObjectResponse, error) {
// 	auther, authErr := r.GetAuther(ctx)
// 	if authErr != nil {
// 		return nil, authErr
// 	}
// 	objectUID, uuidErr := uuid.FromString(uid)
// 	if uuidErr != nil {
// 		return nil, uuidErr
// 	}

// 	response := &graph.GetObjectResponse{}

// 	product, err := r.services.ProductService.GetByUID(ctx, objectUID, auther)
// 	if err == nil {
// 		response.Product = product
// 		return response, nil
// 	}

// 	carton, err := r.services.CartonService.GetByUID(ctx, objectUID, auther)
// 	if err == nil {
// 		response.Carton = carton
// 		return response, nil
// 	}

// 	pallet, err := r.services.PalletService.GetByUID(ctx, objectUID, auther)
// 	if err == nil {
// 		response.Pallet = pallet
// 		return response, nil
// 	}

// 	container, err := r.services.ContainerService.GetByUID(ctx, objectUID, auther)
// 	if err == nil {
// 		response.Container = container
// 		return response, nil
// 	}

// 	return response, nil
// }

// func (r *queryResolver) GetObjects(ctx context.Context, input graph.GetObjectsRequest) (*graph.GetObjectsResponse, error) {
// 	response := &graph.GetObjectsResponse{
// 		Products:   []*models.Product{},
// 		Cartons:    []*models.Carton{},
// 		Pallets:    []*models.Pallet{},
// 		Containers: []*models.Container{},
// 	}

// 	if len(input.ProductUIDs) > 0 {
// 		products, err := r.services.ProductService.GetManyByUIDs(ctx, input.ProductUIDs)
// 		if err != nil {
// 			return nil, err
// 		}
// 		response.Products = products
// 	}

// 	if len(input.CartonUIDs) > 0 {
// 		cartons, err := r.services.CartonService.GetManyByUIDs(ctx, input.CartonUIDs)
// 		if err != nil {
// 			return nil, err
// 		}
// 		response.Cartons = cartons
// 	}

// 	if len(input.PalletUIDs) > 0 {
// 		pallets, err := r.services.PalletService.GetManyByUIDs(ctx, input.PalletUIDs)
// 		if err != nil {
// 			return nil, err
// 		}
// 		response.Pallets = pallets
// 	}

// 	if len(input.ContainerUIDs) > 0 {
// 		containers, err := r.services.ContainerService.GetManyByUIDs(ctx, input.ContainerUIDs)
// 		if err != nil {
// 			return nil, err
// 		}
// 		response.Containers = containers
// 	}

// 	return response, nil
// }

///////////////
// Mutations //
///////////////

// func (r *mutationResolver) RequestToken(ctx context.Context, input *graph.RequestToken) (string, error) {
// 	request := models.LoginRequest{Email: input.Email, Password: input.Password}
// 	auther, err := r.services.AuthService.Login(ctx, &request)
// 	if err != nil {
// 		return "", fmt.Errorf(err.Message)
// 	}

// 	tokenPayload, err := authtoken.Generate(auther)
// 	if err != nil {
// 		return "", fmt.Errorf(err.Message)
// 	}
// 	return tokenPayload.TokenString, nil
// }
