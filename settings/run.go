package settings

func HandleControl(arg string) {
	switch arg {
	case "start":
		StartHandle()

	case "stop":
		StopHandle()

	case "version":
		GetVersion()

	default:
		DefaultHandle()
	}
}
