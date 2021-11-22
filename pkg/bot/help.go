package bot

import "github.com/barklan/cto/pkg/storage"

func getHelpString(data *storage.Data) string {
	return `*Commands:*
/start - start me
/clear - delete all messages from me
/mute - mute for several hours
/unmute - unmute
/status - show current status
/checks - show enabled checks
/log - TBD

` + "```\n/log environment hh[:[m[m[:[s[s]]]]]] [service_name] [flag]\n```"
}
