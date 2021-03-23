package controllerss



func (s *Server) initializeRoutes() {

	v1 := s.Router.Group("/api/v1")
	{
		// Login Route
		v1.POST("/login", s.Login)



		//Users routes
		v1.POST("/users", s.CreateUser)
		// The user of the app have no business getting all the users.
		// v1.GET("/users", s.GetUsers)
		v1.GET("/users/:id", s.GetUser)
		v1.PUT("/users/:id", middlewares.TokenAuthMiddleware(), s.UpdateUser)
		v1.PUT("/avatar/users/:id", middlewares.TokenAuthMiddleware(), s.UpdateAvatar)
		v1.DELETE("/users/:id", middlewares.TokenAuthMiddleware(), s.DeleteUser)

		//product clothes routes
		v1.POST("/productClothes", middlewares.TokenAuthMiddleware(), s.CreateProductClothes)
		v1.GET("/productClothes", s.GetProductClothes)
		v1.GET("/productClothes/:id", s.GetProductClothes)
		v1.PUT("/productClothes/:id", middlewares.TokenAuthMiddleware(), s.UpdateProductClothes)
		v1.DELETE("/productClothes/:id", middlewares.TokenAuthMiddleware(), s.DeleteProductClothes)
		v1.GET("/user_productClothes/:id", s.GetProductClothes)


	}
}
