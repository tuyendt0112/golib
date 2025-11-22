package queue

// Options configures worker behavior for processing jobs.
// These options control retry logic, concurrency, and job lifecycle.
type Options struct {
	// Priority determines job execution order (1-10000).
	// Higher priority jobs are processed first.
	// Default: 0 (no priority)
	Priority uint
	
	// MaxFails is the maximum number of times a job can fail before being moved to dead queue.
	// If set to 1, job goes straight to dead queue on first failure (unless SkipDead is true).
	// Default: 3 (from defaultMaxFails in worker.go)
	MaxFails uint
	
	// SkipDead, when true, prevents failed jobs from being moved to dead queue.
	// Failed jobs are simply discarded after MaxFails attempts.
	// Useful for jobs where retries don't make sense (e.g., invalid data).
	SkipDead bool
	
	// MaxConcurrency limits the number of jobs processed simultaneously by this worker.
	// 0 means no limit (process as many as possible).
	// Useful for rate limiting or resource management.
	MaxConcurrency uint
	
	// MaxTimeout is the maximum time (in seconds) a job can run before being cancelled.
	// 0 means no timeout (job can run indefinitely).
	// Prevents jobs from hanging and consuming resources forever.
	MaxTimeout uint
}

// WithPriority returns an option function to set job priority.
// Higher priority jobs are processed before lower priority ones.
//
// Example:
//   worker := queue.NewWorker("my-queue", queue.WithPriority(100))
func WithPriority(number uint) func(*Options) {
	return func(o *Options) {
		o.Priority = number
	}
}

// WithMaxFails returns an option function to set maximum retry attempts.
// After MaxFails failures, the job is moved to dead queue (unless SkipDead is true).
//
// Example:
//   worker := queue.NewWorker("my-queue", queue.WithMaxFails(5))
func WithMaxFails(number uint) func(*Options) {
	return func(o *Options) {
		o.MaxFails = number
	}
}

// WithSkipDead returns an option function to skip dead queue.
// Failed jobs are discarded instead of being moved to dead queue.
//
// Example:
//   worker := queue.NewWorker("my-queue", queue.WithSkipDead())
func WithSkipDead() func(*Options) {
	return func(o *Options) {
		o.SkipDead = true
	}
}

// WithMaxConcurrency returns an option function to limit concurrent job processing.
// Prevents worker from processing too many jobs at once.
//
// Example:
//   worker := queue.NewWorker("my-queue", queue.WithMaxConcurrency(5))
func WithMaxConcurrency(number uint) func(*Options) {
	return func(o *Options) {
		o.MaxConcurrency = number
	}
}

// WithMaxTimeout returns an option function to set job timeout in seconds.
// Jobs exceeding this timeout are cancelled.
//
// Example:
//   worker := queue.NewWorker("my-queue", queue.WithMaxTimeout(300)) // 5 minutes
func WithMaxTimeout(number uint) func(*Options) {
	return func(o *Options) {
		o.MaxTimeout = number
	}
}