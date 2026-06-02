package section

type (
	Processor struct {
		WebServer ProcessorWebServer `split_world:"true"`
	}

	ProcessorWebServer struct {
		ListenPort uint32 `default:"8080" split_world:"true"`
	}
)
