package chat

import "errors"

var AttemptsExceeded = errors.New("the number of attempts exceeded")
var Close = errors.New("close")
