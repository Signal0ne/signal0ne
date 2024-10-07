package errors

import "fmt"

var ErrConditionNotSatisfied = fmt.Errorf("condition not satisfied")
var ErrAlertAlreadyInactive = fmt.Errorf("alert already inactive")
