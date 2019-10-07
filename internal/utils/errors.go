package utils

import (
	"fmt"

	"github.com/shomali11/slacker"
)

func ReplyWithError(e error, m string, response slacker.ResponseWriter) {
	if e != nil {
		response.ReportError(fmt.Errorf("%s error: %s", m, e.Error()))
	}
}
