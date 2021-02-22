package main

import (
	"context"
	graceful "github.com/nnqq/scr-lib-graceful"
	"github.com/nnqq/scr-org/area"
	"github.com/nnqq/scr-org/config"
	"github.com/nnqq/scr-org/consumer"
	"github.com/nnqq/scr-org/dadata"
	"github.com/nnqq/scr-org/location"
	"github.com/nnqq/scr-org/logger"
	"github.com/nnqq/scr-org/manager"
	"github.com/nnqq/scr-org/metro"
	"github.com/nnqq/scr-org/mongo"
	"github.com/nnqq/scr-org/okved"
	"github.com/nnqq/scr-org/org"
	"github.com/nnqq/scr-org/orgimpl"
	"github.com/nnqq/scr-org/stan"
	pbOrg "github.com/nnqq/scr-proto/codegen/go/org"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"strings"
	"sync"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	stanConn, err := stan.NewConn(cfg.ServiceName, cfg.STAN.ClusterID, cfg.NATS.URL)
	logg.Must(err)

	db, err := mongo.NewConn(ctx, cfg.ServiceName, cfg.MongoDB.URL)
	logg.Must(err)

	orgModel := org.NewModel(db)
	areaModel := area.NewModel(db)
	locationModel := location.NewModel(db)
	managerModel := manager.NewModel(db)
	okvedModel := okved.NewModel(db)
	metroModel := metro.NewModel(db)

	cons := consumer.NewConsumer(
		logg.ZL,
		stanConn,
		cfg.ServiceName,
		dadata.NewClient(strings.Split(cfg.DaData.Tokens, ","), db),
		orgModel,
		areaModel,
		locationModel,
		managerModel,
		okvedModel,
		metroModel,
	)
	logg.Must(cons.Subscribe())

	srv := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
	pbOrg.RegisterOrgServer(srv, orgimpl.NewServer(
		logg.ZL,
		orgModel,
		areaModel,
		locationModel,
		managerModel,
		okvedModel,
		metroModel,
	))

	lis, err := net.Listen("tcp", strings.Join([]string{
		"0.0.0.0",
		cfg.Grpc.Port,
	}, ":"))
	logg.Must(err)

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		graceful.HandleSignals(srv.GracefulStop, cons.GracefulStop)
	}()
	go func() {
		defer wg.Done()
		logg.Must(srv.Serve(lis))
	}()
	go func() {
		defer wg.Done()
		logg.Must(cons.Serve())
	}()
	wg.Wait()
}
