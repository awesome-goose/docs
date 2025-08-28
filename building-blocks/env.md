<!-- // Env
//
// Env is a simple environment variable manager that allows you to get and set environment variables.
// It defines an interface for managing environment variables, which can be implemented by different types of environment managers.
// A file environment manager is provided that loads environment variables from a file in the ".env" format.
//
// TODO:
// - redis env loader
// - sql env loader -->

<!--
// Boot initializes env by loading environment variables from the specified directory.
// If the directory is empty, it defaults to the current working directory.
// It also loads environment variables from the server and sets them.
// The env file is expected to be named ".env" and is in the "env" format.
func (e \*env) Boot() {
viper.AddConfigPath(e.directory)
viper.SetConfigFile(".env")
viper.SetConfigType("env")

    viper.AutomaticEnv()

    envVars := os.Environ()
    for _, env := range envVars {
    	parts := strings.SplitN(env, "=", 2)
    	if len(parts) == 2 {
    		key := parts[0]
    		value := parts[1]
    		viper.Set(key, value)
    	}
    }

    if err := viper.ReadInConfig(); err != nil {
    	log.Println("Error reading env file", err)
    }

}

Env, COnfig key and values must all be a string
string util package provide String() struct with ToBool, ToNumber etc for m=conversion -->
