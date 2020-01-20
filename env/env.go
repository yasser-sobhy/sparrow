package env

type Env struct {
	environment string
}

var env = Env{}

func New() Env { return Env{"development"} }

func Is(e string) bool { return env.environment == e }

func IsDevelopment() bool { return env.environment == "development" }
func IsStaging() bool     { return env.environment == "staging" }
func IsProduction() bool  { return env.environment == "production" }

func SetEnvironment(e string) { env.environment = e }
