package support

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/jmorgan1321/SpaceRep/v1/internal/debug"
)

func OpenFile(filename string) ([]byte, error) {
	debug.Trace()
	defer debug.UnTrace()

	f, err := os.Open(filename)
	if err != nil {
		return nil, LogError("error:", err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, LogError("error:", err)
	}

	return data, nil
}

func ReadData(data []byte) (interface{}, error) {
	debug.Trace()
	defer debug.UnTrace()

	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, LogError("error:", err)
	}
	return v, nil
}
