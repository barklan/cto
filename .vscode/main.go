package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

const start = `{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "up:stack",
            "type": "shell",
            "command": "bash -i run.sh ${input:runsh}",
            "isBackground": true,
        }
    ],
    "inputs": [
        {
            "description": "run.sh command",
            "id": "runsh",
            "type": "pickString",
            "options": [
`

const end = `
            ]
        }
    ]
}
`

func runShCmds() []string {
	content, err := ioutil.ReadFile("../run.sh")
	if err != nil {
		log.Panicln("failed to read run.sh", err)
	}

	r, err := regexp.Compile(`function (\S*) {`)
	if err != nil {
		log.Panicln(err)
	}

	matches := r.FindAllStringSubmatch(string(content), -1)

	cmds := make([]string, 0)
	for _, match := range matches {
		cmd := match[1]
		if !(string(cmd[0]) == "_" || cmd == "help") {
			cmds = append(cmds, match[1])
		}
	}
	return cmds
}

func main() {
	cmds := runShCmds()
	var opts string
	for _, cmd := range cmds {
		opts += fmt.Sprintf(`"%s", `, cmd)
	}
	taskfile := start + opts + end

	err := ioutil.WriteFile("tasks.json", []byte(taskfile), 0777)
	if err != nil {
		panic(err)
	}
}
