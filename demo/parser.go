package demo

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Parser is function which parses Command from array of tokens
type Parser func(tokens []string) (Command, error)

// parsers contain parser functions for all available commands
// key is the name of command
// value is the parser for remaining tokens (without name)
var parsers = map[string]Parser{
	"wait": requireOneFloat(func(f float64) (Command, error) {
		return &Wait{Seconds: f}, nil
	}),
	"hideMessageBox": requireZeroParams(func() (Command, error) {
		return &HideMessageBox{}, nil
	}),
	"setAngle": requireOneInt(func(i int64) (Command, error) {
		return &SetAngle{Angle: int(i)}, nil
	}),
	"shoot": requireOneInt(func(i int64) (Command, error) {
		return &Shoot{Power: int(i)}, nil
	}),
	"waitForFinishTurn": requireZeroParams(func() (Command, error) {
		return &WaitForFinishTurn{}, nil
	}),
	"nextRound": requireZeroParams(func() (Command, error) {
		return &NextRound{}, nil
	}),
	"moveFocus": requireZeroParams(func() (Command, error) {
		return &MoveFocus{}, nil
	}),
	"pressButton": requireZeroParams(func() (Command, error) {
		return &PressButton{}, nil
	}),
	"exit": requireZeroParams(func() (Command, error) {
		return &Exit{}, nil
	}),
}

// ParseScript will read all commands from reader r.
// Commands can be separated by newline character or by `;` (semicolon).
// Empty lines and lines starting with '#' (comments) are ignored.
func ParseScript(r io.Reader) (Script, error) {
	script := Script{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		for _, row := range strings.Split(line, ";") {
			tokens := strings.Fields(row)
			parse := parsers[tokens[0]]
			if parse == nil {
				return nil, fmt.Errorf("Unknown command '%s'", tokens[0])
			}
			command, err := parse(tokens[1:])
			if err != nil {
				return nil, fmt.Errorf("Cannot parse command '%s': %w", tokens[0], err)
			}
			script = append(script, command)
		}
	}
	return script, nil
}

// requireNParams will validate that number of tokens is n before invoking given parser
func requireNParams(n int, parse Parser) Parser {
	return func(tokens []string) (Command, error) {
		if len(tokens) != n {
			return nil, fmt.Errorf("Invalid parameter count %d, expected %d", len(tokens), n)
		}
		return parse(tokens)
	}
}

// requireZeroParams will validate that tokens are empty before invoking given parser
func requireZeroParams(parse func() (Command, error)) Parser {
	return requireNParams(0, func(tokens []string) (Command, error) { return parse() })
}

// requireOneFloat will validate that tokens contain exactly one float number before invoking given parser
func requireOneFloat(parse func(f float64) (Command, error)) Parser {
	return requireNParams(1, func(tokens []string) (Command, error) {
		f, err := strconv.ParseFloat(tokens[0], 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid parameter type: %w", err)
		}
		return parse(f)
	})
}

// requireOneInt will validate that tokens contain exactly one int number before invoking given parser
func requireOneInt(parse func(i int64) (Command, error)) Parser {
	return requireNParams(1, func(tokens []string) (Command, error) {
		i, err := strconv.ParseInt(tokens[0], 0, 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid parameter type: %w", err)
		}
		return parse(i)
	})
}
