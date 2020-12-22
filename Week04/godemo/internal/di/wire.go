// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package di

import (
	"goDemo/internal/dao"
	"goDemo/internal/service"
	"goDemo/internal/server/grpc"
	"goDemo/internal/server/http"

	"github.com/google/wire"
)

//go:generate kratos t wire
func InitApp() (*App, func(), error) {
	panic(wire.Build(dao.Provider, service.Provider, http.New, grpc.New, NewApp))
}
