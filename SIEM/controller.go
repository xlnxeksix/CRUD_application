package SIEM

import (
	models "awesomeProject1/Models"
	"awesomeProject1/SIEM/Model"
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
	var formRule Model.RuleForm

	if err := c.ShouldBindJSON(&formRule); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(formRule)
	// Get the appropriate insight strategy based on the SIEM product type
	selectedStrategy := ctrl.Strategies[formRule.Product]

	ARule := selectedStrategy.InsightAnalysis(formRule)
	fmt.Println(ARule.WildCardInsight)

}
