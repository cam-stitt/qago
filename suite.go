package qago

// Predicate is used to find a selection within the page
type Predicate struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Selector string `yaml:"selector"`
	Text     string `yaml:"text"`
	Multi    bool   `yaml:"multi"`
	First    bool   `yaml:"first"`
}

// KV is a query attribute to assert
type KV struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// Assertion is a validation of the current state
type Assertion struct {
	Name       string `yaml:"name"`
	Text       string `yaml:"text"`
	Query      []KV   `yaml:"query"`
	Attributes []KV   `yaml:"attributes"`
	URL        string `yaml:"url"`
	Visible    bool   `yaml:"visible"`
}

// Action is an event to perform on the selection
type Action struct {
	Name       string      `yaml:"name"`
	Type       string      `yaml:"type"`
	Text       string      `yaml:"text"`
	Wait       string      `yaml:"wait"`
	Assertions []Assertion `yaml:"assertions"`
}

// Step is a step in our test suite
type Step struct {
	Name       string     `yaml:"name"`
	Predicate  *Predicate `yaml:"predicate"`
	Actions    []Action   `yaml:"actions"`
	Wait       string     `yaml:"wait"`
	MaxRetries uint32     `yaml:"max_retries"`
}

// Case is a particular set of steps to be run
type Case struct {
	Name       string      `yaml:"name"`
	Location   string      `yaml:"location"`
	Browser    string      `yaml:"browser"`
	Steps      []Step      `yaml:"steps"`
	Assertions []Assertion `yaml:"assertions"`
}
