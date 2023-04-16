package config

type Logger struct {
	// Log level for zerolog
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	// ConsoleLogEnabled to log using stdout or stderr console
	ConsoleLogEnabled bool `env:"CONSOLE_LOG_ENABLED"`

	// ConsoleLogNoColor to disable color in console log
	ConsoleLogNoColor bool `env:"CONSOLE_LOG_NO_COLOR"`

	// FileLogEnabled to log using file
	FileLogEnabled bool `env:"FILE_LOG_ENABLED"`

	// Directory to log to to when filelogging is enabled
	Directory string `env:"FILE_LOG_DIRECTORY" envDefault:"logs"`

	// Filename is the name of the logfile which will be placed inside the directory
	Filename string `env:"FILE_LOG_NAME" envDefault:"server.log"`

	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int `env:"FILE_LOG_MAX_SIZE" envDefault:"100"`

	// MaxBackups the max number of rolled files to keep
	MaxBackups int `env:"FILE_LOG_MAX_BACKUPS" envDefault:"0"`

	// MaxAge the max age in days to keep a logfile
	MaxAge int `env:"FILE_LOG_MAX_AGE" envDefault:"7"`

	// LocalTime choose the log time format locally or not
	LocalTime bool `env:"FILE_LOG_LOCALTIME"`

	// Compress the log file
	Compress bool `env:"FILE_LOG_COMPRESSED"`
}
