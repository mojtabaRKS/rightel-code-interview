package command

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mojtabaRKS/rightel-code-interview/internal/config"
	"github.com/mojtabaRKS/rightel-code-interview/internal/infra"
	"github.com/mojtabaRKS/rightel-code-interview/internal/repository"
	orderService "github.com/mojtabaRKS/rightel-code-interview/internal/service/order"
)

const expirationBatchSize = 100

type ExpireReservations struct {
	Logger *logrus.Logger
}

func (cmd ExpireReservations) Command(ctx context.Context, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "expire-reservations",
		Short: "expire a batch of pending reservations",
		RunE: func(_ *cobra.Command, _ []string) error {
			return cmd.run(ctx, cfg)
		},
	}
}

func (cmd ExpireReservations) run(ctx context.Context, cfg *config.Config) error {
	postgresClient, err := infra.NewPostgresClient(ctx, cfg.Database.Postgres)
	if err != nil {
		return fmt.Errorf("connect to postgresql: %w", err)
	}

	orderRepo := repository.NewOrderRepository(postgresClient.GetDb())
	orderSvc := orderService.NewOrderService(orderRepo)
	if err := orderSvc.ExpirePending(ctx, expirationBatchSize); err != nil {
		return fmt.Errorf("expire pending reservations: %w", err)
	}

	cmd.Logger.WithContext(ctx).Info("finished expiring pending reservations")
	return nil
}
