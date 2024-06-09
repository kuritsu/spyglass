package commands

import (
	"bytes"
	"encoding/json"
	"os"
)

func DisplayJson(result any) {
	jsonBytes, _ := json.Marshal(result)
	var out bytes.Buffer
	json.Indent(&out, jsonBytes, "", "  ")
	out.WriteTo(os.Stdout)
}
