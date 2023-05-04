package plugins

// simple plugin interface
type Plugin interface {
	API() // URL https://auth.estuary.tech/

	// structure of the input must be defined
	Input(interface{}) error

	// structure of the output must be defined
	Output() interface{} // can be an error
}

// flow
// 1. load all the plugins in this directory
// 2. override pre-existing plugin with the plugins from this directory
