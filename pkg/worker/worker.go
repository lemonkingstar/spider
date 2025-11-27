package worker

type Worker interface {
	// Start starts worker.
	Start() error
	// Stop stops worker.
	Stop() error
}
