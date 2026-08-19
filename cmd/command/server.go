package command

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mojtabaRKS/rightel-code-interview/internal/api"
	orderHandler "github.com/mojtabaRKS/rightel-code-interview/internal/api/handler/order"
	productHandler "github.com/mojtabaRKS/rightel-code-interview/internal/api/handler/product"
	"github.com/mojtabaRKS/rightel-code-interview/internal/config"
	"github.com/mojtabaRKS/rightel-code-interview/internal/infra"
	"github.com/mojtabaRKS/rightel-code-interview/internal/repository"
	orderService "github.com/mojtabaRKS/rightel-code-interview/internal/service/order"
	productService "github.com/mojtabaRKS/rightel-code-interview/internal/service/product"
)

type Server struct {
	Logger *logrus.Logger
}

func (cmd Server) Command(ctx context.Context, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "run Gateway server",
		Run: func(_ *cobra.Command, _ []string) {
			cmd.main(cfg, ctx)
		},
	}
}

func (cmd Server) main(cfg *config.Config, ctx context.Context) {
	postgresClient, err := infra.NewPostgresClient(ctx, cfg.Database.Postgres)
	if err != nil {
		cmd.Logger.WithContext(ctx).Fatal(errors.Wrap(err, "server : failed to connect to postgresql"))
		return
	}

	orderRepo := repository.NewOrderRepository(postgresClient.GetDb())
	orderSvc := orderService.NewOrderService(orderRepo)
	orderH := orderHandler.NewOrderHandler(orderSvc)
	productRepo := repository.NewProductRepository(postgresClient.GetDb())
	productSvc := productService.NewProductService(productRepo)
	productH := productHandler.NewProductHandler(productSvc)

	server := api.New(cfg.AppEnv)
	server.SetupAPIRoutes(orderH, productH)

	if err := server.Serve(ctx, fmt.Sprintf(":%d", cfg.HTTP.Port)); err != nil {
		cmd.Logger.Fatal(err)
	}
}
