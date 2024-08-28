package mio

// 📚 pkg: mio - defines transports for message io

type (
	// Transport provides a reader/writer conduit through which
	// io can take place
	Transport interface {
		// Address returns the location with which this transport interacts
		Address() string

		// Read reads data from the transport at the defined address
		Read() ([]byte, error)

		// Write writes data to the transport at the defined address
		Write(data []byte) error
	}
)
