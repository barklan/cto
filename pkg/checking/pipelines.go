package checking

import (
	"fmt"
	"log"

	"github.com/barklan/cto/pkg/gitlab"
	"github.com/barklan/cto/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

func CheckFailedPipelines(b *tb.Bot, data *storage.Data, args ...interface{}) {
	branch := args[0].(string)

	failedPipelines, err1 := gitlab.GetPipelines(branch, "failed")
	succeededPipelines, err2 := gitlab.GetPipelines(branch, "success")
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
		data.CSend(fmt.Sprintf("@barklan. Pipeline %s failed on %s.", fmt.Sprint(latestFailID), branch))
	}
}
