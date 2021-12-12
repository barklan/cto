package bot

// func registerOnTextHanler(b *tb.Bot, data *storage.Data) {
// 	b.Handle(tb.OnText, func(m *tb.Message) {
// 		project, ok := VerifySender(data, m)
// 		if !ok {
// 			return
// 		}

// 		if !strings.Contains(m.Text, "=") {
// 			return
// 		}
// 		lowerText := strings.ToLower(m.Text)
// 		split := strings.SplitN(lowerText, "=", 2)
// 		varName, valueToSet := split[0], split[1]

// 		kv := reflect.ValueOf(storage.KVars)
// 		for i := 0; i < kv.NumField(); i++ {
// 			if varName == kv.Field(i).String() {
// 				data.SetVar(project, varName, valueToSet, -1)
// 			}
// 		}

// 		simFloat, err := strconv.ParseFloat(simStr, 64)
// 		if err != nil {
// 			data.CSend("Failed to parse float.")
// 		}

// 		data.Config.Internal.Log.SimilarityThreshold = simFloat
// 		data.CSend("New threshold is set.")
// 	})
// }
