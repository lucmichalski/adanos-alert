package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/internal/repository/impl"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

var Version string
var GitCommit string

func main() {
	app := glacier.Create(fmt.Sprintf("%s (%s)", Version, GitCommit))
	app.WithHttpServer(":19999")

	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "mongo_uri",
		Usage:  "Mongodb connection uri",
		EnvVar: "ADANOS_MONGODB_URI",
		Value:  "mongodb://localhost:27017",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "mongo_db",
		Usage:  "Mongodb database name",
		EnvVar: "ADANOS_MONGODB_DB",
		Value:  "adanos",
	}))

	app.Singleton(func(c *cli.Context) *configs.Config {
		return &configs.Config{
			MongoURI: c.String("mongo_uri"),
			MongoDB:  c.String("mongo_db"),
		}
	})

	app.Singleton(func(conf *configs.Config) *mongo.Database {
		conn, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(conf.MongoURI))
		if err != nil {
			log.Errorf("connect to mongodb failed: %s", err)
			return nil
		}

		return conn.Database(conf.MongoDB)
	})

	app.Singleton(impl.NewSequenceRepo)
	app.Singleton(impl.NewKVRepo)
	app.Singleton(impl.NewMessageRepo)
	app.Singleton(impl.NewMessageGroupRepo)

	if err := app.Run(os.Args); err != nil {
		log.Errorf("exit with error: %s", err)
	}
}
