# env.MustGet

`MustGet` is a small utility function designed to safely read environment variables in Go applications.

It ensures that important configuration values (database URL, JWT secrets, third-party API keys, etc.) **must exist**, otherwise the application will not start.  
This prevents unexpected runtime errors caused by missing configuration.

---

## 🔧 How It Works

```go
value, ok := os.LookupEnv(key)


## Example (optional value):

port := os.Getenv("SERVER_PORT")

```
