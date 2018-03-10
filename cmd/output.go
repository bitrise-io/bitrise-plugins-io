package cmd

import (
	"encoding/json"
	"fmt"
)

func printOutput(data []byte, pretty bool) {
	var output string
	if pretty {
		var outModel map[string]interface{}
		if err := json.Unmarshal(data, &outModel); err != nil {
			printErrorOutput(err.Error(), pretty)
			return
		}

		out, err := json.MarshalIndent(outModel["data"], "", " ")
		if err != nil {
			printErrorOutput(err.Error(), pretty)
			return
		}
		output = string(out)
	} else {
		output = string(data)
	}
	fmt.Printf(output)
}

func printErrorOutput(message string, pretty bool) {
	if pretty {
		fmt.Println(message)
	} else {
		fmt.Printf(`{"error":"%s"}`, message)
	}
}
