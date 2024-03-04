package container

import (
	"boilerplate/config"
	"boilerplate/internal/app"
	"boilerplate/internal/filesystem"
	"boilerplate/internal/infra/database"
	"boilerplate/internal/infra/http/controllers"
	"boilerplate/internal/infra/http/middlewares"
	"log"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

type Container struct {
	Middlewares
	Services
	Controllers
}

type Middlewares struct {
	AuthMw func(http.Handler) http.Handler
}

type Services struct {
	app.AuthService
	app.UserService
	app.FarmService
	app.CategoryService
	app.OfferService
	app.OrderService
	app.OrderItemsService
	app.ImageModelService
	app.AddressService
}

type Controllers struct {
	controllers.AuthController
	controllers.UserController
	controllers.FarmController
	controllers.CategoryController
	controllers.OfferController
	controllers.OrderController
	controllers.OrderItemController
	controllers.ImageModelController
	controllers.AddressController
}

func New(conf config.Configuration) Container {
	tknAuth := jwtauth.New("HS256", []byte(conf.JwtSecret), nil)
	sess := getDbSess(conf)

	userRepository := database.NewUserRepository(sess)
	sessionRepository := database.NewSessRepository(sess)
	offerRepository := database.NewOfferRepository(sess)
	farmRepository := database.NewFarmRepository(sess, offerRepository)
	orderItemRepository := database.NewOrderItemRepository(sess, offerRepository, farmRepository)
	orderRepository := database.NewOrderRepository(sess, orderItemRepository)
	ImageRepository := database.NewImageModelRepository(sess)
	addressRepository := database.NewAddressRepository(sess)

	userService := app.NewUserService(userRepository)
	authService := app.NewAuthService(sessionRepository, userService, conf, tknAuth)
	farmService := app.NewFarmService(farmRepository, offerRepository, orderRepository)
	catService := app.NewCategoryService()
	imageStorageService := filesystem.NewImageStorageService(conf.FileStorageLocation)
	imageService := app.NewImageModelService(ImageRepository, imageStorageService)
	offerService := app.NewOfferService(offerRepository, imageStorageService, imageService)
	orderService := app.NewOrderService(orderRepository, orderItemRepository)
	orderItemService := app.NewOrderItemsService(orderItemRepository, orderRepository)
	addressService := app.NewAddressService(addressRepository)

	authController := controllers.NewAuthController(authService, userService)
	userController := controllers.NewUserController(userService)
	farmController := controllers.NewFarmController(farmService)
	categoryController := controllers.NewCategoryController(catService)
	offerController := controllers.NewOfferController(offerService, farmService, imageService)
	orderController := controllers.NewOrderController(orderService, orderItemService)
	orderItemController := controllers.NewOrderItemController(orderItemService)
	imageController := controllers.NewImageModelController(imageService)
	addressController := controllers.NewAddressController(addressService)

	authMiddleware := middlewares.AuthMiddleware(tknAuth, authService, userService)

	return Container{
		Middlewares: Middlewares{
			AuthMw: authMiddleware,
		},
		Services: Services{
			authService,
			userService,
			farmService,
			catService,
			offerService,
			orderService,
			orderItemService,
			imageService,
			addressService,
		},
		Controllers: Controllers{
			authController,
			userController,
			farmController,
			categoryController,
			offerController,
			orderController,
			orderItemController,
			imageController,
			addressController,
		},
	}
}

func getDbSess(conf config.Configuration) db.Session {
	sess, err := postgresql.Open(
		postgresql.ConnectionURL{
			User:     conf.DatabaseUser,
			Host:     conf.DatabaseHost,
			Password: conf.DatabasePassword,
			Database: conf.DatabaseName,
		})
	//sess, err := sqlite.Open(
	//	sqlite.ConnectionURL{
	//		Database: conf.DatabasePath,
	//	})
	if err != nil {
		log.Fatalf("Unable to create new DB session: %q\n", err)
	}
	return sess
}
