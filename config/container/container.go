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
	app.AddressService
}

type Controllers struct {
	controllers.AuthController
	controllers.UserController
	controllers.FarmController
	controllers.CategoryController
	controllers.OfferController
	controllers.AddressController
}

func New(conf config.Configuration) Container {
	tknAuth := jwtauth.New("HS256", []byte(conf.JwtSecret), nil)
	sess := getDbSess(conf)

	userRepository := database.NewUserRepository(sess)
	sessionRepository := database.NewSessRepository(sess)
	farmRepository := database.NewFarmRepository(sess)
	offerRepository := database.NewOfferRepository(sess, farmRepository)
	addressRepository := database.NewAddressepository(sess)

	userService := app.NewUserService(userRepository)
	authService := app.NewAuthService(sessionRepository, userService, conf, tknAuth)
	farmService := app.NewFarmService(farmRepository)
	catService := app.NewCategoryService()
	imageService := filesystem.NewImageStorageService(conf.FileStorageLocation)
	offerService := app.NewOfferService(offerRepository, imageService)
	addressService := app.NewAddressService(addressRepository)

	authController := controllers.NewAuthController(authService, userService)
	userController := controllers.NewUserController(userService)
	farmController := controllers.NewFarmController(farmService, userService)
	categoryController := controllers.NewCategoryController(catService)
	offerController := controllers.NewOfferController(offerService, farmService)
	addressController := controllers.NewAddressController(addressService, userService)

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
			addressService,
		},
		Controllers: Controllers{
			authController,
			userController,
			farmController,
			categoryController,
			offerController,
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
