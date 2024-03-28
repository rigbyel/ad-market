package validate

import (
	"regexp"

	"github.com/rigbyel/ad-market/internal/models/constraints"
)

// validates user login
func ValidateLogin(login string) []string {
	if login == "" {
		return []string{"login is required"}
	}

	errs := []string{}

	// check if login has a valid size
	if len(login) < constraints.LoginMinLen {
		errs = append(errs, "login is too short")
	}

	if len(login) > constraints.LoginMaxLen {
		errs = append(errs, "login is too long")
	}

	// check if login consists only from alphanumeric characters
	alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(login)
	if !alphanumeric {
		errs = append(errs, "login should contain only alphanumeric characters")
	}

	return errs
}
