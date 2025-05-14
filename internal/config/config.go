package config

import (
    "flag"
    "fmt"
    "io/ioutil"
    "time"

    "gopkg.in/yaml.v3"
)

type Config struct {
    Host          string        `yaml:"host"`
    Port          int           `yaml:"port"`
    StatsInterval time.Duration // parsed duration
    rawInterval   string        `yaml:"statsInterval"`
}

func Load() (*Config, error) {
    // allow overriding the path if needed
    path := flag.String("config", "config.yaml", "path to YAML config file")
    flag.Parse()

    // read YAML
    data, err := ioutil.ReadFile(*path)
    if err != nil {
        return nil, fmt.Errorf("reading config file %s: %w", *path, err)
    }

    // unmarshal
    var c Config
    if err := yaml.Unmarshal(data, &c); err != nil {
        return nil, fmt.Errorf("parsing config: %w", err)
    }

    // apply defaults
    if c.Host == "" {
        c.Host = "0.0.0.0"
    }
    if c.Port == 0 {
        c.Port = 8080
    }
    if c.rawInterval == "" {
        c.rawInterval = "1h"
    }

    // parse the duration string
    d, err := time.ParseDuration(c.rawInterval)
    if err != nil {
        return nil, fmt.Errorf("invalid statsInterval %q: %w", c.rawInterval, err)
    }
    c.StatsInterval = d

    return &c, nil
}
