package actions

type Result interface {
	Error() error
}

type Value interface {
	Int() int
}
