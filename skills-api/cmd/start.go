/*
Copyright © 2022 Himanshu Shekhar <himanshu.kiit@gmail.com>
Code ownership is with Himanshu Shekhar. Use without modifications.
*/
package cmd

import (
	"fmt"

	"github.com/imhshekhar47/go-api-core/config"
	"github.com/imhshekhar47/go-api-core/logger"
	"github.com/imhshekhar47/go-api-core/server"
	"github.com/imhshekhar47/go-api-core/utils"
	"github.com/imhshekhar47/hs-taskmaster/skills-api/api"
	"github.com/imhshekhar47/hs-taskmaster/skills-api/db"
	"github.com/imhshekhar47/hs-taskmaster/skills-api/pb"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "embed"

	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

var (
	startLogger = logger.GetLogger("cmd/start")
)

//go:embed skill-api.swagger.json
var swaggerDoc string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Strat API service",
	Long:  `Start Task Management API service`,
	Run:   runStartCmd,
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().Uint16VarP(&argPortRest, "rest-port", "r", 0, "Port for REST APIs")
	startCmd.Flags().Uint16VarP(&argPortGrpc, "grpc-port", "g", 50051, "Port for gRPC APIs")
}

func runRest(port uint16, basePath string, grpcAddress string) error {
	restServer := server.NewRestServer(port, basePath, grpcAddress, swaggerDoc)

	pb.RegisterSkillServiceHandlerFromEndpoint(
		restServer.GetContext(),
		restServer.GetMultiplex(),
		restServer.GetGrpcAddress(),
		restServer.GetDialoptions(),
	)

	return restServer.Run()
}

func runGrpc(port uint16, dbObj *gorm.DB) error {
	startLogger.Debugf("entry: runGrpc(%d)", port)
	grpcServer := server.NewGrpcServer(port)

	if dbObj != nil {
		dbObj.AutoMigrate(&pb.SkillORM{})
		pb.RegisterSkillServiceServer(grpcServer.Get(), api.NewServer(db.NewSkillRepository(dbObj)))
	} else {
		pb.RegisterSkillServiceServer(grpcServer.Get(), api.NewServer(db.GetRepositoryLocal()))
	}

	startLogger.Debugf("exit: runGrpc()")
	return grpcServer.Run()
}

func loadDatasource() (*gorm.DB, error) {
	datasource := &config.GetApplicationConfig().Datasource

	if err := datasource.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid datasource")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		datasource.Host,
		datasource.Username,
		datasource.Password,
		datasource.Database,
		datasource.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, sqlErr := db.DB()
	if sqlErr != nil {
		return nil, sqlErr
	}

	pingErr := sqlDB.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return db, nil
}

func runStartCmd(cmd *cobra.Command, args []string) {
	startLogger.Debugln("entry: runStartCmd()")

	startLogger.Tracef("Config: %s", config.GetApplicationConfig().Json())
	startLogger.Tracef("Config: %s", config.GetApplicationConfig().Yaml())

	gormDB, err := loadDatasource()
	if err != nil {
		startLogger.Warnf("failed to load datasource %v", err)
	} else {
		startLogger.Info("database connection established")
		sqlDB, sqlErr := gormDB.DB()
		if sqlErr != nil {
			panic(err)
		}
		defer sqlDB.Close()
	}

	grpcAddr := fmt.Sprintf("0.0.0.0:%d", argPortGrpc)
	// rest
	if argPortRest > 5000 {
		go func() {
			startLogger.Infoln("launcing corouting for rest server")
			err := runRest(argPortRest, utils.GetEnvOrElse("SERVER_BASE_PATH", ""), grpcAddr)
			if err != nil {
				startLogger.Errorln("could not start rest server", err)
			}
		}()
	}

	// grpc
	err = runGrpc(argPortGrpc, gormDB)
	if err != nil {
		startLogger.Errorf("exit: runStartCmd(), Errror: %s", err)
		panic(err)
	}

	startLogger.Debugln("end: runStartCmd()")
}
