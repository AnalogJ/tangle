package middleware

import (
	"fmt"
	"github.com/analogj/tangle/webapp/backend/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/sirupsen/logrus"
)

func DatabaseMiddleware(appConfig config.Interface, globalLogger logrus.FieldLogger) gin.HandlerFunc {
	//var database *gorm.DB
	fmt.Printf("Trying to connect to database stored: %s\n", appConfig.GetString("web.database.location"))

	useConsoleLogger := func(level neo4j.LogLevel) func(config *neo4j.Config) {
		return func(config *neo4j.Config) {
			config.Log = neo4j.ConsoleLogger(level)
		}
	}

	neo4jLevel := neo4j.ERROR
	if appConfig.GetString("log.level") == "DEBUG" {
		neo4jLevel = neo4j.DEBUG
	}

	driver, err := neo4j.NewDriver(appConfig.GetString("web.database.uri"),
		neo4j.BasicAuth(
			appConfig.GetString("web.database.username"),
			appConfig.GetString("web.database.password"),
			appConfig.GetString("web.database.realm"),
		),
		useConsoleLogger(neo4jLevel),
	)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database! - %v", err))
	}
	//defer driver.Close()

	//session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	//defer session.Close()

	//TODO: detrmine where we can call defer database.Close()
	return func(c *gin.Context) {
		c.Set("NEO4J", driver)
		c.Next()
	}
}
