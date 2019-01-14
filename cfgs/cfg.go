package cfgs

import "flag"

type Cfgs struct {
	MySQL struct {
		TcpAddress string
		UserName   string
		Password   string
		DbName     string
		LogAddress string
	}
	Log struct {
		LogAddress string
	}
}

var (
	MODE *string
	//Conf the whole configuration structure
	Conf *Cfgs
)

func init() {
	//define the flag which stands for this application's run mode
	MODE = flag.String("m", "online", "the mode of this application run at")
	Conf = new(Cfgs)
}
