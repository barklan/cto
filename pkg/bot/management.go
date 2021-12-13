package bot

import (
	"fmt"
	"time"

	"github.com/barklan/cto/pkg/storage"
	"github.com/barklan/cto/pkg/storage/vars"
	tb "gopkg.in/tucnak/telebot.v2"
)

func bossChat(data *storage.Data, m *tb.Message) bool {
	return m.Chat.ID == data.Config.Internal.TG.BossChatID
}

func registerManagementHandlers(data *storage.Data, b *tb.Bot) {
	b.Handle("/remove", func(m *tb.Message) {
		if !bossChat(data, m) {
			return
		}
		// TODO magic string
		data.SetVar(storage.Internal, vars.ChatState, "removeCalled", 3*time.Minute)
		data.CSend("Send me a project to remove.")
	})
}

func bossText(data *storage.Data, m *tb.Message) {
	state := string(data.GetVar(storage.Internal, vars.ChatState))
	switch state {
	case "removeCalled":
		if _, ok := data.Config.P[m.Text]; ok {
			storage.DeleteProject(data, m.Text)
			data.CSend("Project removed.")
		} else {
			data.CSend(fmt.Sprintf("Project %q not found.", m.Text))
		}
		return
	}
}

func projectText(data *storage.Data, project string, m *tb.Message) {
	// if !strings.Contains(m.Text, "=") {
	// 	return
	// }
	// lowerText := strings.ToLower(m.Text)
	// split := strings.SplitN(lowerText, "=", 2)
	// varName, valueToSet := split[0], split[1]

	// kv := reflect.ValueOf(storage.KVars)
	// for i := 0; i < kv.NumField(); i++ {
	// 	if varName == kv.Field(i).String() {
	// 		data.SetVar(project, varName, valueToSet, -1)
	// 	}
	// }

	// simFloat, err := strconv.ParseFloat(simStr, 64)
	// if err != nil {
	// 	data.CSend("Failed to parse float.")
	// }

	// data.Config.Internal.Log.SimilarityThreshold = simFloat
	// data.CSend("New threshold is set.")
}

func registerOnTextHanler(data *storage.Data, b *tb.Bot) {
	b.Handle(tb.OnText, func(m *tb.Message) {
		if bossChat(data, m) {
			bossText(data, m)
			return
		}

		if project, ok := VerifySender(data, m); ok {
			projectText(data, project, m)
		}
	})
}
