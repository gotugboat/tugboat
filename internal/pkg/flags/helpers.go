package flags

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
	"tugboat/internal/pkg/slices"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func addFlag(cmd *cobra.Command, flag *Flag) {
	if flag == nil || flag.Name == "" {
		return
	}

	var flags *pflag.FlagSet
	if flag.Persistent {
		flags = cmd.PersistentFlags()
	} else {
		flags = cmd.Flags()
	}

	// Evaluate what type of flag should being added
	switch v := flag.Value.(type) {
	case int:
		flags.IntP(flag.Name, flag.Shorthand, v, flag.Usage)
	case string:
		flags.StringP(flag.Name, flag.Shorthand, v, flag.Usage)
	case []string:
		flags.StringSliceP(flag.Name, flag.Shorthand, v, flag.Usage)
	case bool:
		flags.BoolP(flag.Name, flag.Shorthand, v, flag.Usage)
	case time.Duration:
		flags.DurationP(flag.Name, flag.Shorthand, v, flag.Usage)
	}

	if flag.Deprecated {
		flags.MarkHidden(flag.Name)
	}
}

func bind(cmd *cobra.Command, flag *Flag) error {
	if flag == nil {
		return nil
	} else if flag.Name == "" {
		// This flag is only available in the config file
		viper.SetDefault(flag.ConfigName, flag.Value)
		return nil
	}

	if flag.Persistent {
		if err := viper.BindPFlag(flag.ConfigName, cmd.PersistentFlags().Lookup(flag.Name)); err != nil {
			return err
		}
	} else {
		// fmt.Println("Config Name:", flag.ConfigName, " Flag Name:", cmd.Flags().Lookup(flag.Name))
		if err := viper.BindPFlag(flag.ConfigName, cmd.Flags().Lookup(flag.Name)); err != nil {
			return err
		}
	}

	if strings.Contains(flag.ConfigName, "-") {
		var str string
		replace_chars := []string{"-", "."}

		for _, char := range replace_chars {
			str = strings.ReplaceAll(flag.ConfigName, char, "_")
		}
		str = strings.ToUpper(str)

		viper.BindEnv(flag.ConfigName, str)
	}

	// Bind the yaml configs to an env var (i.e log.format -> LOG_FORMAT)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	return nil
}

func getString(flag *Flag) string {
	if flag == nil {
		return ""
	}
	return viper.GetString(flag.ConfigName)
}

func getStringSlice(flag *Flag) []string {
	if flag == nil {
		return nil
	}
	v := viper.GetStringSlice(flag.ConfigName)

	// Separate env values containing a ','
	switch {
	case len(v) == 0: // no strings
		return nil
	case len(v) == 1 && strings.Contains(v[0], ","): // , separated string
		v = strings.Split(v[0], ",")
	}
	return v
}

func getBool(flag *Flag) bool {
	if flag == nil {
		return false
	}
	return viper.GetBool(flag.ConfigName)
}

type sanitizedInput struct {
	RawInput string
	Value    string
	IsCmd    bool
}

// evaluateValue will attempt to safely evaluate a variable that has been set to a bash expression
// but not run any string as a command unless it is prefixed with `$`. Valid formats are $VALUE, ${VALUE}, $(cat VALUE).
// allowed bash commands to run are defined in the sanitizeInput function
func evaluateValue(value string) (string, error) {
	// if the string does not start with a $ don't evaluate it
	if !strings.HasPrefix(value, "$") {
		return value, nil
	}

	value, err := executeBashCmd(value)
	if err != nil {
		return "", err
	}

	return value, nil
}

// executeBashCmd sanitizes a bash command and executes it.
func executeBashCmd(input string) (string, error) {
	sanitizedInput, err := sanitizeInput(input)
	if err != nil {
		return "", err
	}

	if !sanitizedInput.IsCmd {
		return sanitizedInput.Value, nil
	}

	cmd := exec.Command("bash", "-c", input)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return trim(string(output)), nil
}

func sanitizeInput(input string) (*sanitizedInput, error) {
	var sanitizedValue string
	allowedCommands := []string{"echo", "cat", "git"}
	disallowedCharacters := []string{";", "|", "&"}

	// Use a regular expression to extract Bash expressions from the input string
	r := regexp.MustCompile(`\$\(.+?\)|\${.+?}|\$[A-Z0-9_]+`)
	match := r.FindString(input)

	// Execute the Bash expressions and replace them with their output

	// Check if the match is a $(VALUE) pattern
	if strings.HasPrefix(match, "$(") && strings.HasSuffix(match, ")") {
		// extract the expression to execute if it contains allowed commands
		sanitizedValue = strings.Replace(input, match, match[2:len(match)-1], -1)
	} else {
		// Check if the match is a ${VALUE} pattern
		if strings.HasPrefix(match, "${") && strings.HasSuffix(match, "}") {
			cmd := exec.Command("bash", "-c", "echo "+match)
			output, err := cmd.Output()
			if err != nil {
				return nil, err
			}
			sanitizedValue = strings.Replace(input, match, string(output), -1)
			return &sanitizedInput{RawInput: input, Value: trim(sanitizedValue), IsCmd: false}, nil
		}

		// Check if the match is a $VALUE pattern
		if strings.HasPrefix(match, "$") {
			// Get the value of the environment variable
			value, ok := os.LookupEnv(match[1:])
			if !ok {
				return nil, fmt.Errorf("undefined environment variable: %s", match[1:])
			}
			sanitizedValue = strings.Replace(input, match, value, -1)
			return &sanitizedInput{RawInput: input, Value: trim(sanitizedValue), IsCmd: false}, nil
		}
	}

	// Split the input into words
	words := strings.Fields(sanitizedValue)

	// Check if the first word is an allowed command
	if !slices.Contains(allowedCommands, words[0]) {
		return nil, fmt.Errorf("disallowed command: %s", words[0])
	}

	// Check if the input contains any disallowed characters
	for _, char := range disallowedCharacters {
		if strings.Contains(sanitizedValue, char) {
			return nil, fmt.Errorf("disallowed character: %s", char)
		}
	}

	return &sanitizedInput{RawInput: input, Value: trim(sanitizedValue), IsCmd: true}, nil
}

func trim(input string) string {
	return strings.TrimSuffix(input, "\n")
}
