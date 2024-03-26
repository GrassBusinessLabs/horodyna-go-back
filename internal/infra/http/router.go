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
		AllowedOrigins:   []string{"*"},
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
				OfferRouter(apiRouter, cont.OfferController, cont.OfferService, cont.ImageModelService)
				OrderRouter(apiRouter, cont.OrderController, cont.OrderService, cont.FarmService)
				OrderItemRoute(apiRouter, cont.OrderItemController, cont.OrderService, cont.OrderItemsService)
				ImageRouter(apiRouter, cont.ImageModelController, cont.ImageModelService)
				AddressRouter(apiRouter, cont.AddressController, cont.AddressService)
				InvoiceRouter(apiRouter, cont.InvoiceController, cont.InvoiceService)
				MonobankRouter(apiRouter, cont.MonobankController)

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

func OrderItemRoute(r chi.Router, oc controllers.OrderItemController, os app.OrderService, o app.OrderItemsService) {
	pathObjectMiddleware := middlewares.PathObject("orderId", controllers.OrderKey, os)
	pathObjectItemMiddleware := middlewares.PathObject("orderItemId", controllers.OrderItemKey, o)
	isOwnerMiddleware := middlewares.IsOwnerMiddleware[domain.Order](controllers.OrderKey)
	isOwnerMiddlewareOrderItems := middlewares.IsOwnerMiddleware[domain.OrderItem](controllers.OrderItemKey)

	r.Route("/order-items", func(apiRouter chi.Router) {
		apiRouter.With(pathObjectMiddleware, isOwnerMiddleware).Post(
			"/{orderId}",
			oc.AddItem(),
		)
		apiRouter.With(pathObjectItemMiddleware, isOwnerMiddlewareOrderItems).Put(
			"/{orderItemId}",
			oc.Update(),
		)
		apiRouter.With(pathObjectItemMiddleware, isOwnerMiddlewareOrderItems).Delete(
			"/{orderItemId}",
			oc.Delete(),
		)
	})
}

func OrderRouter(r chi.Router, oc controllers.OrderController, os app.OrderService, fs app.FarmService) {
	pathObjectMiddleware := middlewares.PathObject("orderId", controllers.OrderKey, os)
	isOwnerMiddleware := middlewares.IsOwnerMiddleware[domain.Order](controllers.OrderKey)
	farmPathObjectMiddleware := middlewares.PathObject("farmId", controllers.FarmKey, fs)
	farmIsOwnerMiddleweare := middlewares.IsOwnerMiddleware[domain.Farm](controllers.FarmKey)
	r.Route("/orders", func(apiRouter chi.Router) {
		apiRouter.With(pathObjectMiddleware, isOwnerMiddleware).Put(
			"/receiver-status/{orderId}",
			oc.SetOrderStatusAsReceiver(),
		)
		apiRouter.With(pathObjectMiddleware, farmPathObjectMiddleware, farmIsOwnerMiddleweare).Put(
			"/farmer-status/{farmId}/{orderId}",
			oc.SetOrderStatusAsFarmer(),
		)
		apiRouter.Get(
			"/farmer-percentage",
			oc.GetFarmerOrdersPercentage(),
		)
		apiRouter.With(pathObjectMiddleware).Get(
			"/split/{orderId}",
			oc.SplitOrderByFarms(),
		)
		apiRouter.With(pathObjectMiddleware, isOwnerMiddleware).Post(
			"/split/{orderId}/{farmId}",
			oc.SubmitSplitedOrder(),
		)
		apiRouter.With(pathObjectMiddleware, isOwnerMiddleware).Delete(
			"/split/{orderId}/{farmId}",
			oc.DeleteSplitedOrder(),
		)
		apiRouter.Get(
			"/by-farmer",
			oc.FindByFarmUserId(),
		)
		apiRouter.With(pathObjectMiddleware).Get(
			"/{orderId}",
			oc.FindById(),
		)
		apiRouter.With(pathObjectMiddleware, isOwnerMiddleware).Put(
			"/{orderId}",
			oc.Update(),
		)
		apiRouter.With(pathObjectMiddleware, isOwnerMiddleware).Delete(
			"/{orderId}",
			oc.Delete(),
		)
		apiRouter.Get(
			"/",
			oc.FindAllByUserId(),
		)
		apiRouter.Post(
			"/",
			oc.Save(),
		)
	})
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
		apiRouter.Post(
			"/get-by-coords",
			uc.FindAllByCoords(),
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

func OfferRouter(r chi.Router, oc controllers.OfferController, os app.OfferService, is app.ImageModelService) {

	pathObjectMiddleware := middlewares.PathObject("offerId", controllers.OfferKey, os)
	imagePathObjectMiddleware := middlewares.PathObject("imageId", controllers.ImageKey, is)
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
		apiRouter.Get(
			"/by-farmid/{farmId}",
			oc.FindByFarmId(),
		)
		apiRouter.With(pathObjectMiddleware, isOwnerMiddleware).Post(
			"/additional-image/{offerId}",
			oc.AddAdditionalImage(),
		)
		apiRouter.With(pathObjectMiddleware, imagePathObjectMiddleware, isOwnerMiddleware).Delete(
			"/additional-image/{offerId}/{imageId}",
			oc.DeleteAdditionalImage(),
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

func AddressRouter(r chi.Router, ac controllers.AddressController, as app.AddressService) {
	pathObjectMiddleware := middlewares.PathObject("addressId", controllers.AddressKey, as)
	isOwnerMiddleware := middlewares.IsOwnerMiddleware[domain.Address](controllers.AddressKey)

	r.Route("/address", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/",
			ac.Save(),
		)
		apiRouter.With(pathObjectMiddleware).Get(
			"/{addressId}",
			ac.FindById(),
		)
		apiRouter.With(pathObjectMiddleware, isOwnerMiddleware).Put(
			"/{addressId}",
			ac.Update(),
		)
		apiRouter.With(pathObjectMiddleware, isOwnerMiddleware).Delete(
			"/{addressId}",
			ac.Delete(),
		)
		apiRouter.Get(
			"/by-user/{userId}",
			ac.FindByUserId(),
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
		apiRouter.Post(
			"/login-email",
			ac.LoginWithEmail(),
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
			"/phone-number",
			uc.SetPhoneNumber(),
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
		apiRouter.Get(
			"/",
			ic.FindAll(),
		)
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

func InvoiceRouter(r chi.Router, ic controllers.InvoiceController, is app.InvoiceService) {
	pathObjectMiddleware := middlewares.PathObject("invoiceId", controllers.InvoiceKey, is)
	r.Route("/invoice", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/",
			ic.FindAll(),
		)
		apiRouter.With(pathObjectMiddleware).Get(
			"/{invoiceId}",
			ic.FindById(),
		)
		apiRouter.Get(
			"/last-day",
			ic.FindAllUpdatedWithinOneDay(),
		)
	})

}

func MonobankRouter(r chi.Router, mc controllers.MonobankController) {
	r.Route("/monobank", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/",
			mc.CreateInvoice(),
		)
		apiRouter.Get(
			"/{invoiceId}",
			mc.GetInvoiceData(),
		)
		apiRouter.Post(
			"/cancel",
			mc.CreateInvoice(),
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
