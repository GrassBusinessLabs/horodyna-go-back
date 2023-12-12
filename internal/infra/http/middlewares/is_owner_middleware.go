package middlewares

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/controllers"
	"errors"
	"net/http"
)

func IsOwnerMiddleware[domainType controllers.Userable](key controllers.CtxKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user := ctx.Value(controllers.GetUserKey()).(domain.User)
			obj := controllers.GetPathValFromCtx[domainType](ctx, key)

			if obj.GetUserId() != user.Id {
				err := errors.New("you have no access to this object")
				controllers.Forbidden(w, err)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
