package main

import (
	"context"
	graceful "github.com/leaq-ru/lib-graceful"
	"github.com/leaq-ru/org/area"
	"github.com/leaq-ru/org/config"
	"github.com/leaq-ru/org/consumer"
	"github.com/leaq-ru/org/dadata"
	"github.com/leaq-ru/org/location"
	"github.com/leaq-ru/org/logger"
	"github.com/leaq-ru/org/manager"
	"github.com/leaq-ru/org/metro"
	"github.com/leaq-ru/org/mongo"
	"github.com/leaq-ru/org/okved"
	"github.com/leaq-ru/org/org"
	"github.com/leaq-ru/org/orgimpl"
	"github.com/leaq-ru/org/stan"
	pbOrg "github.com/leaq-ru/proto/codegen/go/org"
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
