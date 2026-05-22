package command

import (
	_ "embed"
)

const (
	CommandIndex      = 0
	FirstTargetIndex  = 0
	SecondTargetIndex = 1
	ThirdTargetIndex  = 2
	FourthTargetIndex = 3
	FifthTargetIndex  = 4
	SixthTargetIndex  = 5
)

//go:embed Default.txt
var DefaultManual string

//go:embed Simple.txt
var SimpleManual string

//go:embed Advance.txt
var AdvanceManual string
