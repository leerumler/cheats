package ggconf

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
)

// Xinfos holds information about the current X connection state.
type Xinfos struct {
	XUtil        *xgbutil.XUtil
	Root, Active *xproto.Window
}
