package ReadCfg

//

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type CfgType struct {
	HostPort  string        `json:"HostPort"`
	SleepTime time.Duration `json:"SleepTime"`
}

func ReadCfg(fn string) (cfg CfgType) {
	cfg.HostPort = "localhost:8080"
	cfg.SleepTime = time.Duration(5)
	buf, err := ioutil.ReadFile(fn)
	if err == nil {
		err = json.Unmarshal(buf, &cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse %s for configuratin, error=%s\n", fn, err)
			os.Exit(1)
		}
	}
	return
}

/* vim: set noai ts=4 sw=4: */
