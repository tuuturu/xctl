package sdk

import "io"

type helmSDKClient struct {
	debugOut io.Writer
	debug    bool
}

type NewHelmClientOpts struct {
	DebugOut io.Writer
	Debug    bool
}
