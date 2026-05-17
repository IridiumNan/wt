package model

type (
	WTMethod uint8
	WTType   uint8
)

const (
	WTRead    WTMethod = iota
	WTInstall WTMethod = iota
	WTWrite   WTMethod = iota
)

const (
	WTServer WTType = iota
	WTClient WTType = iota
)

type Auth struct {
	WtMethod WTMethod
	Token    string
	ErrMsg   string
}
