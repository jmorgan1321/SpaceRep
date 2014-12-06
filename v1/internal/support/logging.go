package support

import (
	"fmt"
	"github.com/jmorgan1321/SpaceRep/v1/internal/debug"
	"log"
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
