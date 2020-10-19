package main

import (
        "os"
        "gopkg.in/yaml.v2"
)
 
type Config struct {
        Database   Database   `yaml:"database"`
}
 
type Database struct {
        Influxdb Influxdb `yaml:"influxdb"`
}
 
type Influxdb struct {
        Host            string `yaml:"host"`
        Port            string `yaml:"port"`
        Username        string `yaml:"username"`
        Password        string `yaml:"password"`
        Database        string `yaml:"database"`
        RetentionPolicy string `yaml:"retention_policy"`
}
 
func GetConfig(configfile string) (*Config, error) {
        // buka berkas konfigurasi
        file, err := os.Open(configfile)
        if err != nil {
                return nil, err
        }
        defer file.Close()
 
        config := &Config{}
 
        d := yaml.NewDecoder(file)
 
        if err := d.Decode(&config); err != nil {
                return nil, err
        }
 
        return config, nil
}