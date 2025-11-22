package queue

// Options is a struct that holds the options for the worker
type Options struct {
	Priority       uint // Priority from 1 to 10000
	MaxFails       uint // 1: send straight to dead (unless SkipDead)
	SkipDead       bool // If true, don't send failed jobs to the dead queue when retries are exhausted.
	MaxConcurrency uint // Max number of jobs to keep in flight (default is 0, meaning no max)
	MaxTimeout     uint // Max time in seconds to wait for a job to finish (default is 0, meaning no timeout)
}

// WithPriority is a function that sets the priority of the worker
func WithPriority(number uint) func(*Options) {

	return func(o *Options) {
		o.Priority = number
	}
}

// WithMaxFails is a function that sets the maximum number of fails for the worker
func WithMaxFails(number uint) func(*Options) {

	return func(o *Options) {
		o.MaxFails = number
	}
}

// WithSkipDead is a function that sets the worker to skip dead jobs
func WithSkipDead() func(*Options) {

	return func(o *Options) {
		o.SkipDead = true
	}
}

// WithMaxConcurrency is a function that sets the maximum concurrency of the worker
func WithMaxConcurrency(number uint) func(*Options) {

	return func(o *Options) {
		o.MaxConcurrency = number
	}
}

// WithMaxTimeout is a function that sets the maximum timeout of the worker
func WithMaxTimeout(number uint) func(*Options) {
	return func(o *Options) {
		o.MaxTimeout = number
	}
}