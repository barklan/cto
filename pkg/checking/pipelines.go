package checking

import (
	"fmt"
	"log"

	"github.com/barklan/cto/pkg/gitlab"
	"github.com/barklan/cto/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

func CheckFailedPipelines(b *tb.Bot, data *storage.Data, args ...interface{}) {
	projectName := args[0].(string)
	branch := args[1].(string)

	gitlabProjectId := fmt.Sprint(data.Config.P[projectName].Checks.GitLab.ProjectID)
	gitlabToken := data.Config.P[projectName].Checks.GitLab.APIToken
	failedPipelines, err1 := gitlab.GetPipelines(gitlabProjectId, gitlabToken, branch, "failed")
	succeededPipelines, err2 := gitlab.GetPipelines(gitlabProjectId, gitlabToken, branch, "success")
	if err1 != nil || err2 != nil {
		log.Println("failed to get gitlab pipelines")
		return
	}

	var latestFailID int64
	for _, failedPipeline := range failedPipelines {
		if failedPipeline.ID > latestFailID {
			latestFailID = failedPipeline.ID
		}
	}

	var latestSuccessID int64
	for _, successPipeline := range succeededPipelines {
		if successPipeline.ID > latestSuccessID {
			latestSuccessID = successPipeline.ID
		}
	}

	if latestFailID > latestSuccessID {
		data.PSend(projectName, fmt.Sprintf("Pipeline %s failed on %s.", fmt.Sprint(latestFailID), branch))
	}
}
