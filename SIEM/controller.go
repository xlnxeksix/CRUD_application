package SIEM

import (
	models "awesomeProject1/Models"
	"awesomeProject1/SIEM/Model"
	"awesomeProject1/SIEM/Strategies"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Import statements...

type Controller struct {
	Repo       SIEMRepository
	Strategies map[string]Strategies.Insight
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
		Strategies: strategies,
	}
}

func (ctrl *Controller) GetRuleContent(c *gin.Context) {
	var formRule Model.RuleForm

	if err := c.ShouldBindJSON(&formRule); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate the duration until the next scheduled time
	now := time.Now()
	nextExecutionTime := time.Date(now.Year(), now.Month(), now.Day(), formRule.ScheduledTime.Hour(), formRule.ScheduledTime.Minute(), formRule.ScheduledTime.Second(), 0, now.Location())

	if now.After(nextExecutionTime) {
		nextExecutionTime = nextExecutionTime.Add(24 * time.Hour)
	}

	durationUntilNextExecution := nextExecutionTime.Sub(now)

	// Schedule the task to run at the next execution time
	timer := time.NewTimer(durationUntilNextExecution)

	// Run the task when the timer fires and pass formRule as an argument
	go func(formRule Model.RuleForm) {
		fmt.Println("hi")
		<-timer.C

		// Get the appropriate insight strategy based on the SIEM product type
		selectedStrategy := ctrl.Strategies[formRule.Product]

		InsightIDs := selectedStrategy.InsightAnalysis(&formRule)
		flattenedRule := &Model.FlattenedRule{
			RuleForm:      formRule,
			FlattenedRule: formRule.RuleContent, // Assuming RuleContent is the flattened rule
		}
		ctrl.Repo.InsertInsight(flattenedRule, InsightIDs)
		fmt.Println("Added")

		// Schedule the task to run again at the same time the next day
		go ctrl.scheduleTaskAtUserTime(formRule)
	}(formRule)
}

func (ctrl *Controller) scheduleTaskAtUserTime(formRule Model.RuleForm) {
	for {
		fmt.Println("hi")
		now := time.Now()
		nextExecutionTime := time.Date(now.Year(), now.Month(), now.Day(), formRule.ScheduledTime.Hour(), formRule.ScheduledTime.Minute(), formRule.ScheduledTime.Second(), 0, now.Location()).Add(24 * time.Hour)

		durationUntilNextExecution := nextExecutionTime.Sub(now)

		time.Sleep(durationUntilNextExecution)

		selectedStrategy := ctrl.Strategies[formRule.Product]

		InsightIDs := selectedStrategy.InsightAnalysis(&formRule)
		flattenedRule := &Model.FlattenedRule{
			RuleForm:      formRule,
			FlattenedRule: formRule.RuleContent, // Assuming RuleContent is the flattened rule
		}
		ctrl.Repo.InsertInsight(flattenedRule, InsightIDs)
		fmt.Println("Added")
		fmt.Println("Added")
	}
}
