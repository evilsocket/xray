package grabbers

import xray "github.com/evilsocket/xray"

func Init() {
	xray.SetupGrabbers(
		[]xray.Grabber{
			&HTTPGrabber{},
			&DNSGrabber{},
			NewLineGrabber("smtp", []int{25, 587}),
			NewLineGrabber("ftp", []int{21}),
			NewLineGrabber("ssh", []int{22, 222, 2222}),
			NewLineGrabber("pop", []int{110}),
			NewLineGrabber("irc", []int{6667}),
		})
}
