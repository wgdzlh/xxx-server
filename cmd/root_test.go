package cmd

import (
	"os"
	"testing"
	"time"

	"xxx-server/domain/entity"
	"xxx-server/infrastructure/config"

	json "github.com/json-iterator/go"
)

func TestOutputWsPng(t *testing.T) {
	config.C.Server.RunLocal = true
	config.C.Server.CmdRoot = ".."
	SetupConfig()
	config.C.Cron.DisableRecover = true
	go serverRun(nil, nil)
	// app.SchedulerSvr.OutputWsPng()
}

func TestArray(t *testing.T) {
	var a = [29]int16{}
	for i := range a {
		a[i] = int16(i + 1)
	}
	t.Logf("%v", a)
	b := [4]int16(a[:])
	b[2] = 13
	t.Logf("%v", b)
	t.Logf("%v", a)
}

func TestMapJson(t *testing.T) {
	a := `{
		"MaxWind": 11.75,
		"WindScale": 6,
		"Alert_Level": "blue",
		"Alert_Info": "预计未来24小时内可能受大风影响，平均风力可达6级以上，请注意防范。",
		"ValidHour": 24
	}`
	m := map[string]entity.AnyJson{}
	err := json.Unmarshal([]byte(a), &m)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(json.MarshalToString(m))
}

func TestChdir(t *testing.T) {
	mainWd, _ := os.Getwd()
	t.Log("mainWd", mainWd)
	go func() {
		t.Log(os.Getwd())
		os.Chdir("../test")
		t.Log(os.Getwd())
	}()
	time.Sleep(time.Second)
	mainWd, _ = os.Getwd()
	t.Log("mainWd", mainWd)
}

func TestJson(t *testing.T) {
	a := `{"adcode": "150000", "lat_lon": [88, 120, 45, 50], "subs": [1, 2, 3], "district": "内蒙古自治区"}`
	c := `{"lat_lon": null, "district": ""}`
	var b entity.AlertDistConfig
	err := json.Unmarshal([]byte(a), &b)
	t.Log(err, b)
	err = json.Unmarshal([]byte(c), &b)
	t.Log(err, b)
}
