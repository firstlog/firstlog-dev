package initialize

// ParseConfig
type Main struct {
	Input  Input  `yaml:"Input"`
	Output Output `yaml:"Output"`
}

type Input struct {
	Tasks  []Task  `yaml:"Tasks"`
}

type Task struct {
	Recursive 	  bool	  `yaml:"Recursive"`
	Directory	  string  `yaml:"Directory"`
	Ignore        string  `yaml:"Ignore"`
	Match		  string  `yaml:"Match"`
}

type Output struct {
	Elasticsearch Elasticsearch `yaml:"Elasticsearch"`
}

type Elasticsearch struct {
	Index    string    `yaml:"Index"`
	Hosts    []string  `yaml:"Hosts"`
	Version  string    `yaml:"Version"`
	Shards   string    `yaml:"Shards"`
	Replicas string    `yaml:"Replicas"`
	Detail   Detail    `yaml:"Detail"`
}

type Detail struct {
	Enable   bool   `yaml:"Enable"`
	Regex    string `yaml:"Regex"`
	Template string `yaml:"Template"`
}
