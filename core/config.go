package core

import (
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)
// Config represents the global configuration
type Config struct {
	CACertificate   []byte    `xml:"CACert"`
	CAPrivateKey    []byte    `xml:"CAPvt"`
	Project					*Project
}

// History represents a history of projects
type History struct {
	H []*Project `xml:"ProjectsHistory"`
}

// Project represents the basic information of a project
type Project struct {
	Title string
	Path  string
}

var globalSettingsFileName = "global_settings.xml"

var historyFileName = "history.xml"

// LoadGlobalSettings loads the global settings
func LoadGlobalSettings(path string) *Config {

	// if path doesn't exists, create it
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0700)
	}

	var cfg *Config
	xmlSettingsPath := filepath.Join(path, globalSettingsFileName)
	xmlSettingsFile, err := os.Open(xmlSettingsPath)
	defer xmlSettingsFile.Close()
	if err != nil {
		// create the file and put in a new fresh default settings
		cfg = initGlobalSettings()
		saveGlobalSettings(cfg, xmlSettingsPath)
		return cfg
	}

	byteValue, err := ioutil.ReadAll(xmlSettingsFile)
	if err != nil {
		fmt.Println(err)
	}

	err = xml.Unmarshal(byteValue, &cfg)
	if err != nil {
		cfg = initGlobalSettings()
		saveGlobalSettings(cfg, xmlSettingsPath)
	}

	// TODO: add a method that checks the configuration just loaded

	return cfg

}

// LoadHistory loads the projects history
func LoadHistory(path string) *History {

	historyPath := filepath.Join(path, historyFileName)
	historyFile, err := os.Open(historyPath)
	defer historyFile.Close()
	if err != nil {
		// no history
		return &History{}
	}
	byteValue, err := ioutil.ReadAll(historyFile)
	// TODO: handle error
	if err != nil {
		fmt.Println(err)
	}
	var history *History
	err = xml.Unmarshal(byteValue, &history)
	if err != nil {
		return &History{}
	}
	return history
}

func initGlobalSettings() *Config {
	// generate a new CA
	// TODO: handle error generated by CreateCA
	rawPvt, rawCA, _ := CreateCA()
	pemPvt := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: rawPvt})
	pemCA := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rawCA})
	cfg := &Config{
		CACertificate: pemCA,
		CAPrivateKey:  pemPvt,
	}
	return cfg
}

func saveGlobalSettings(cfg *Config, path string) error {
	xmlSettings, _ := xml.MarshalIndent(cfg, "", " ")
	return ioutil.WriteFile(path, xmlSettings, 0700)
}
