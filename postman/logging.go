package postman

import (
	"io"
	"io/ioutil"
	"log"
)

func init() {
	// by default throw it all away
	log.SetOutput(ioutil.Discard)
}

func EnableLogging(writer io.Writer) {
	log.SetOutput(writer)
}
