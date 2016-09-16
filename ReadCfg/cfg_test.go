package ReadCfg

//
// Test of reading config file, dependent on ./cfg.json
//

import "testing"

func Test_Hash512(t *testing.T) {

	cfg := ReadCfg("./cfg.json")

	if cfg.HostPort != "localhost:8123" {
		t.Errorf("Error reading ./cfg.json")
	}
	if cfg.SleepTime != 1 {
		t.Errorf("Error reading ./cfg.json")
	}

}

/* vim: set noai ts=4 sw=4: */
