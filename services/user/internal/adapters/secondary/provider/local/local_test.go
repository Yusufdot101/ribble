package local

import "testing"

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		isValid bool
	}{
		{
			name:    "vaild email",
			email:   "example@gmail.com",
			isValid: true,
		},
		{
			name:  "invaild email",
			email: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid := isValidEmail(test.email)
			if isValid != test.isValid {
				t.Errorf("expected isValidEmail=%v, got isValidEmail=%v", test.isValid, isValid)
			}
		})
	}
}
