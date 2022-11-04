package customErrors

type Error interface {
	Error() string
}

type ExpiredToken struct {
}

type EmailNotValidated struct {

}

func (e ExpiredToken) Error() string {
	return "Token was Expired"
}

func (e EmailNotValidated) Error() string {
	return "Email was not validated"
}

type InvalidToken struct {
	E error
}

func (e InvalidToken) Error() string {
	return e.E.Error()
}

