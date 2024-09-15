package validation

type UserValidationConstants struct {
	Username struct {
		MinLength int
		MaxLength int
		Regex     string
	}
	DisplayName struct {
		MinLength int
		MaxLength int
		Regex     string
	}
	Bio struct {
		MinLength int
		MaxLength int
	}
}

var UserValidationConstantsInstance = UserValidationConstants{
	Username: struct {
		MinLength int
		MaxLength int
		Regex     string
	}{
		MinLength: 3,
		MaxLength: 42,
		Regex:     "^[a-z0-9_]{3,42}$",
	},
	DisplayName: struct {
		MinLength int
		MaxLength int
		Regex     string
	}{
		MinLength: 3,
		MaxLength: 42,
		Regex:     "^[a-zA-Z0-9_ ]{3,42}$",
	},
	Bio: struct {
		MinLength int
		MaxLength int
	}{
		MinLength: 0,
		MaxLength: 200,
	},
}
