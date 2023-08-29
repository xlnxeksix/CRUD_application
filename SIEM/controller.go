package SIEM

import (
	models "awesomeProject1/Models"
	"awesomeProject1/SIEM/Strategies"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Controller struct {
	Repo       SIEMRepository
	Strategies map[string]Strategies.Insight // Add the strategies map
}

func NewSIEMController(repo SIEMRepository) *Controller {
	strategies := map[string]Strategies.Insight{
		"splunk":       &Strategies.SplunkStrategy{},
		"splunk_cloud": &Strategies.SplunkCloudStrategy{},
		"qradar":       &Strategies.QradarStrategy{},
		"sentinel":     &Strategies.SentinelStrategy{},
	}
	return &Controller{
		Repo:       repo,
		Strategies: strategies, // Initialize the strategies map
	}
}

func (ctrl *Controller) GetRuleContent(c *gin.Context) {
	var rule rule

	if err := c.ShouldBindJSON(&rule); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the appropriate insight strategy based on the SIEM type
	selectedStrategy := ctrl.Strategies[rule.SIEM]

	// Calculate the shipping price using the selected strategy
	Insight := selectedStrategy.InsightAnalysis(rule.RuleContent)

	fmt.Println(Insight)

}
