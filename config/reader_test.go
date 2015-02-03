package config

import "testing"

func TestNewConfig(t *testing.T) {
	blah, err := NewConfig("../config.json.dist")

	if err != nil {
		t.Errorf("knackered, %s was meant to be nil!", err.Error())
	}

	if blah.Filepath != "../config.json.dist" {
		t.Errorf("filepath %s does not match expected", blah.Filepath)
	}

	if blah.Config.Port != "8080" {
		t.Errorf("%s does not match expected default!", blah.Config.Port)
	}

	if blah.Config.Workers != 10 {
		t.Errorf("expected %d workers by default", blah.Config.Workers)
	}
}

func TestNewConfigIsFukt(t *testing.T) {
	conf, err := NewConfig("some/bollocks")

	if conf != nil {
		t.Errorf("expect config to be nil with a bad path")
	}

	if err == nil {
		t.Errorf("expect error to be populated with a bad path")
	}
}
