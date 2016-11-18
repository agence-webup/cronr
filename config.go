package main

// ConfigFile represents a YAML config file
type ConfigFile struct {
	Jobs []Job
}

// Job represents a unit to schedule
type Job struct {
	Cron    string   `yaml:"cron"`
	Command []string `yaml:"command"`
	// User       uint32   `yaml:"user"`
	// Group      uint32   `yaml:"group"`
	WorkingDir string `yaml:"working_dir"`
}
