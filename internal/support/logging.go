package support

import (
	"fmt"
	"log"

	"github.com/jmorgan1321/SpaceRep/internal/debug"
)

var indent = "\t\t"

func LogError(msg string, err error) error {
	log.Println(debug.IndentationLevel, indent, "error:", err)
	return err
}

func LogFatal(msg string, v ...interface{}) {
	s := fmt.Sprintf("%s"+indent+msg+"\n", append([]interface{}{debug.IndentationLevel}, v...)...)
	panic(s)
}

func Log(msg string, v ...interface{}) {
	log.Printf("%s"+indent+msg+"\n", append([]interface{}{debug.IndentationLevel}, v...)...)
}
