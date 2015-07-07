package main

import (
	"github.com/mikeflynn/go-alexa/skillserver"
)

var Applications = map[string]skillserver.EchoApplication{
	"/echo/test": skillserver.EchoApplication{ // Route
		AppID:   "amzn1.echo-sdk-ams.app.872bfbc9-005e-47f3-a02a-c8c657d4e0f2", // Echo App ID
		Handler: EchoHelloWorld,                                                // Handler Func
	},
	"/echo/jeopardy": skillserver.EchoApplication{
		AppID:   "amzn1.echo-sdk-ams.app.872bfbc9-005e-47f3-a02a-c8c657d4e0f2",
		Handler: EchoJeopardy,
	},
}

func main() {
	skillserver.Run(Applications, "3000")
}
