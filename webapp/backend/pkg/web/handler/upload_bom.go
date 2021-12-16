package handler

import (
	"github.com/analogj/tangle/webapp/backend/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/sirupsen/logrus"
)

func UploadBom(c *gin.Context) {

	neo4jDriver := c.MustGet("NEO4J").(neo4j.Driver)
	logger := c.MustGet("LOGGER").(logrus.FieldLogger)
	appConfig := c.MustGet("CONFIG").(config.Interface)

	session := neo4jDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

}
