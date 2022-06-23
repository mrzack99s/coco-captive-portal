package config

import (
	"context"

	"github.com/mrzack99s/coco-captive-portal/types"
	"github.com/rs/zerolog"
)

var (
	Config                       types.ConfigType
	NetLog                       *zerolog.Logger
	AppLog                       *zerolog.Logger
	PROD_MODE                    bool
	URL_CAPTIVE_PORTAL_DETECTION = []string{
		"http://www.gstatic.com/generate_204",
		"http://clients3.google.com/generate_204",
		"http://www.apple.com/library/test/success",
		"http://connectivitycheck.android.com/generate_204",
		"http://connectivitycheck.gstatic.com/generate_204",
		"http://www.msftncsi.com/ncsi.txt",
		"http://www.appleiphonecell.com",
		"http://captive.apple.com",
		"http://captive.roku.com/ok",
		"http://detectportal.firefox.com/success.txt",
		"http://www.msftconnecttest.com/connecttest.txt",
		"http://fireoscaptiveportal.com/generate_204",
		"http://connectivitycheck.cbg-app.huawei.com/generate_204",
		"http://connect.rom.miui.com/generate_204",
		"http://freetimecaptiveportal.com/generate_204",
		"http://gateway.zscalerthree.net/generate_204",
		"http://gateway.zscloud.net/generate_204",
		"http://g.cn/generate_204",
		"http://play.googleapis.com/generate_204",
		"http://speedtest-global.spatialbuzz.net/generate_204",
		"http://tabletcaptiveportal.com/generate204",
		"http://www.google.cn/generate_204",
		"http://edge.microsoft.com/captiveportal/generate_204",
	}
	DETECTOR_CONTEXT             context.Context
	DETECTOR_CONTEXT_CANCEL_FUNC context.CancelFunc
	AUTH_CONTEXT                 context.Context
	AUTH_CONTEXT_CANCEL_FUNC     context.CancelFunc
	NET_WATCHER_CONTEXT          context.Context
	NET_WATCHER_CANCEL_FUNC      context.CancelFunc
)

func init() {
	DETECTOR_CONTEXT, DETECTOR_CONTEXT_CANCEL_FUNC = context.WithCancel(context.Background())
	AUTH_CONTEXT, AUTH_CONTEXT_CANCEL_FUNC = context.WithCancel(context.Background())
	NET_WATCHER_CONTEXT, NET_WATCHER_CANCEL_FUNC = context.WithCancel(context.Background())
}
