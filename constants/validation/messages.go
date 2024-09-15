package validation

type MessageValidationConstants struct {
	Content struct {
		MinLength int
		MaxLength int
	}
}

var MessageValidationConstantsInstance = MessageValidationConstants{
	Content: struct {
		MinLength int
		MaxLength int
	}{
		MinLength: 1,
		MaxLength: 5000,
	},
}
