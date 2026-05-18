package command

import (
	_ "embed"
)

//go:embed Default.txt
var DefaultManual string

//go:embed Simple.txt
var SimpleManual string

//go:embed Advance.txt
var AdvanceManual string
