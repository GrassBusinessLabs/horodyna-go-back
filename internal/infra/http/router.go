package http

import (
	"boilerplate/config"
	"boilerplate/config/container"
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/controllers"
	"boilerplate/internal/infra/http/middlewares"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Router(cont container.Container) http.Handler {

	router := chi.NewRouter()

	router.Use(middleware.RedirectSlashes, middleware.Logger, cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Route("/api", func(apiRouter chi.Router) {
		// Health
		apiRouter.Route("/ping", func(healthRouter chi.Router) {
			healthRouter.Get("/", PingHandler())
			healthRouter.Handle("/*", NotFoundJSON())
		})

		apiRouter.Route("/v1", func(apiRouter chi.Router) {
			// Public routes
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Route("/auth", func(apiRouter chi.Router) {
					AuthRouter(apiRouter, cont.AuthController, cont.AuthMw)
				})
				CategoryRouter(apiRouter, cont.CategoryController)
			})

			// Protected routes
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Use(cont.AuthMw)

				UserRouter(apiRouter, cont.UserController)
				FarmRouter(apiRouter, cont.FarmController, cont.FarmService)
				OfferRouter(apiRouter, cont.OfferController, cont.OfferService)
				ImageRouter(apiRouter, cont.ImageModelController, cont.ImageModelService)

				apiRouter.Handle("/*", NotFoundJSON())
			})
		})
	})

	router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, config.GetConfiguration().FileStorageLocation))
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
		fs.ServeHTTP(w, r)
	})

	return router
}

func CategoryRouter(r chi.Router, categoryController controllers.CategoryController) {
	r.Route("/categories", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/",
			categoryController.FindAll(),
		)
	})
}

func FarmRouter(r chi.Router, uc controllers.FarmController, fs app.FarmService) {
	pathObjectMiddleware := middlewares.PathObject("farmId", controllers.FarmKey, fs)
	isOwnerMiddleware := middlewares.IsOwnerMiddleware[domain.Farm](controllers.FarmKey)

	r.Route("/farms", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/",
			uc.ListView(),
		)
		apiRouter.With(pathObjectMiddleware).Get(
			"/{farmId}",
			uc.FindById(),
		)
		apiRouter.Post(
			"/",
			uc.Save(),
		)
		apiRouter.With(pathObjectMiddleware, isOwnerMiddleware).Put(
			"/{farmId}",
			uc.Update(),
		)
		apiRouter.With(pathObjectMiddleware, isOwnerMiddleware).Delete(
			"/{farmId}",
			uc.Delete(),
		)
	})
}

func OfferRouter(r chi.Router, oc controllers.OfferController, os app.OfferService) {

	pathObjectMiddleware := middlewares.PathObject("offerId", controllers.OfferKey, os)
	isOwnerMiddleware := middlewares.IsOwnerMiddleware[domain.Offer](controllers.OfferKey)

	r.Route("/offers", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/",
			oc.Save(),
		)
		apiRouter.Get(
			"/",
			oc.ListView(),
		)
		apiRouter.With(pathObjectMiddleware).Get(
			"/{offerId}",
			oc.FindById(),
		)
		apiRouter.With(pathObjectMiddleware, isOwnerMiddleware).Delete(
			"/{offerId}",
			oc.Delete(),
		)
		apiRouter.With(pathObjectMiddleware, isOwnerMiddleware).Put(
			"/{offerId}",
			oc.Update(),
		)
	})
}

func AuthRouter(r chi.Router, ac controllers.AuthController, amw func(http.Handler) http.Handler) {
	r.Route("/", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/register",
			ac.Register(),
		)
		apiRouter.Post(
			"/login",
			ac.Login(),
		)
		apiRouter.With(amw).Post(
			"/change-pwd",
			ac.ChangePassword(),
		)
		apiRouter.With(amw).Post(
			"/logout",
			ac.Logout(),
		)
	})
}

func UserRouter(r chi.Router, uc controllers.UserController) {
	r.Route("/users", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/",
			uc.FindMe(),
		)
		apiRouter.Put(
			"/",
			uc.Update(),
		)
		apiRouter.Delete(
			"/",
			uc.Delete(),
		)
	})
}

func ImageRouter(r chi.Router, ic controllers.ImageModelController, is app.ImageModelService) {
	pathObjectMiddleware := middlewares.PathObject("imageId", controllers.ImageKey, is)
	r.Route("/images", func(apiRouter chi.Router) {
		apiRouter.With(pathObjectMiddleware).Get(
			"/{imageId}",
			ic.FindById(),
		)
		apiRouter.Post(
			"/",
			ic.Save(),
		)

		apiRouter.With(pathObjectMiddleware).Delete(
			"/{imageId}",
			ic.Delete(),
		)
	})

}

func NotFoundJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode("Resource Not Found")
		if err != nil {
			fmt.Printf("writing response: %s", err)
		}
	}
}

func PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode("Ok")
		if err != nil {
			fmt.Printf("writing response: %s", err)
		}
	}
}
