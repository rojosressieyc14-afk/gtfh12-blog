package utils

import (
	"fmt"
	"unicode"
)

type PasswordPolicy struct {
	MinLength      int
	MaxLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireDigit   bool
	RequireSpecial bool
}

var DefaultPasswordPolicy = PasswordPolicy{
	MinLength:      8,
	MaxLength:      64,
	RequireUpper:   true,
	RequireLower:   true,
	RequireDigit:   true,
	RequireSpecial: false,
}

func ValidatePassword(password string, policy PasswordPolicy) error {
	if len(password) < policy.MinLength {
		return fmt.Errorf("密码长度不能少于 %d 位", policy.MinLength)
	}
	if len(password) > policy.MaxLength {
		return fmt.Errorf("密码长度不能超过 %d 位", policy.MaxLength)
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if policy.RequireUpper && !hasUpper {
		return fmt.Errorf("密码必须包含至少一个大写字母")
	}
	if policy.RequireLower && !hasLower {
		return fmt.Errorf("密码必须包含至少一个小写字母")
	}
	if policy.RequireDigit && !hasDigit {
		return fmt.Errorf("密码必须包含至少一个数字")
	}
	if policy.RequireSpecial && !hasSpecial {
		return fmt.Errorf("密码必须包含至少一个特殊字符")
	}

	return nil
}
