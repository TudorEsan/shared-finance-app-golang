package customErrors

type Error interface {
	Error() string
}

type ExpiredToken struct {
}

func (e ExpiredToken) Error() string {
	return "Token was Expired"
}

type InvalidToken struct {
}

func (e InvalidToken) Error() string {
	return "Token is Invalid"
}

func GetJsonError(e Error) string {
	return e.Error()
}
