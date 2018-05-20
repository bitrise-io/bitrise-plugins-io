package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/bitrise-io/go-utils/log"
	"github.com/pkg/errors"
)

// PrettyOutput ...
type PrettyOutput interface {
	Pretty() string
}

func printOutputWithPrettyFormatter(data []byte, pretty bool, prettyFormatter PrettyOutput) error {
	var output string
	if pretty {
		if err := json.Unmarshal(data, prettyFormatter); err != nil {
			return errors.WithStack(err)
		}
		output = prettyFormatter.Pretty()
	} else {
		output = string(data)
	}
	fmt.Println(output)
	return nil
}

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
	fmt.Println(output)
}

func printErrorOutput(message string, pretty bool) {
	if pretty {
		log.Errorf(message)
	} else {
		fmt.Printf(`{"error":"%s"}`, message)
	}
}
