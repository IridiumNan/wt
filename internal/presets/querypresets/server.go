package querypresets

import (
	"gitee.com/cai-zixiang_hainan/wt/internal/model"
)

var ServerHostPortQuery = model.Query{
	Head:    ">>set the host and port in the server<<",
	Example: "0.0.0.0:12212",
	Default: "0.0.0.0:12212",
}

var ServerReadTimeoutQuery = model.Query{
	Head:    ">>set the read timeout for server<<",
	Example: "20s",
	Default: "20s",
}

var ServerInstallTimeoutQuety = model.Query{
	Head:    ">>set the install timeout for server<<",
	Example: "3h",
	Default: "3h",
}

var ServerWriteTimeoutQuery = model.Query{
	Head:    ">>set the write timeout for server<<",
	Example: "3h",
	Default: "3h",
}

var ServerReadTokenQuery = model.Query{
	Head:    ">>set the read tokens for server<<",
	Example: "xxxxx",
	Default: "no token set as default",
}

var ServerInstallTokenQuery = model.Query{
	Head:    ">>set the install tokens for server<<",
	Example: "xxxxx",
	Default: "no token set as default",
}

var ServerWriteTokenQuery = model.Query{
	Head:    ">>set the write tokens for server",
	Example: "xxxxx",
	Default: "no token set as default",
}
