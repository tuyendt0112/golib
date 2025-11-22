package env

import "os"

// MustGet retrieves the value of an environment variable.
//
// It works like this:
//   - If the environment variable with the given "key" exists → return its value.
//   - If it does NOT exist → panic immediately and stop the program.
//     This is useful for required variables (DB_URL, JWT_SECRET, AWS_KEY...),
//     because it prevents the server from starting with incorrect configuration.
func MustGet(key string) string {
	// LookupEnv returns the value and a boolean indicating if the variable is set
	variable, ok := os.LookupEnv(key)
	// If the variable is NOT set, panic immediately
	if !ok {
		panic("environment variable " + key + " not found")
	}
	// Otherwise, return the value
	return variable
}
