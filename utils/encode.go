// It takes a data structure and an output stream, encodes the data structure into JSON format, and writes it to the output
// stream
package utils

import (
	"encoding/json"
	"io"
)

// Encoding the data into JSON format and writing it to the output stream.
func PrettyEncode(data interface{}, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return nil
}
