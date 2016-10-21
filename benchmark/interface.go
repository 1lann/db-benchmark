package benchmark

var Wrappers = make(map[string]Wrapper)

// Document represents a key value document
type Document struct {
	ID    string
	Value string
}

// ConnectOpts represents the options used to connect to a generic database
type ConnectOpts struct {
	DB       string
	Table    string
	Username string
	Password string
	Host     string
}

// Wrapper represents a database query wrapper
type Wrapper interface {
	Connect(opts ConnectOpts)

	Get(id string) Document
	GetAll() []Document

	Update(doc Document)
	Put(doc Document)

	Clear()
}

// RegisterWrapper registers a wrapper
func RegisterWrapper(name string, wrapper Wrapper) {
	Wrappers[name] = wrapper
}
