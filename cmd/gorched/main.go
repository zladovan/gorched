package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/zladovan/gorched"
	"github.com/zladovan/gorched/demo"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	app := &cli.App{
		Name:  "gorched",
		Usage: "console game inspired by Scorched Earth written in GO",
		UsageText: `
		For the best game experience maximize your terminal before start.
		
		Start with --seed to play the same sequence of rounds again.
		After finishing game you will see initial seed and seed for the last round.`,
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:        "seed",
				Usage:       "Integer `NUMBER` used as seed for random generations",
				DefaultText: "current time",
				Aliases:     []string{"s"},
			},
			&cli.IntFlag{
				Name:        "width",
				Usage:       "Width of the game world in `NUMBER` of console cells",
				DefaultText: "actual terminal width",
			},
			&cli.IntFlag{
				Name:        "height",
				Usage:       "Height of the game world in `NUMBER` of console cells",
				DefaultText: "actual terminal height",
			},
			&cli.IntFlag{
				Name:  "fps",
				Usage: "Screen framerate, use lower values to reduce system resources usage",
				Value: 40,
			},
			&cli.BoolFlag{
				Name:  "ascii-only",
				Usage: "Use only ASCII characters to draw graphics",
			},
			&cli.BoolFlag{
				Name:  "low-color",
				Usage: "Use only 8 colors to draw graphics",
			},
			&cli.BoolFlag{
				Name:  "browser",
				Usage: "Use this flag when starting from emulated terminal in browser",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Turn on debug mode",
			},
			&cli.StringFlag{
				Name:  "demo",
				Usage: "Play demo script from given `FILE` right after game start",
			},
		},
		HideHelpCommand: true,
		Action:          run,
	}

	// run application
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	// init seed
	seed := c.Int64("seed")
	if seed == 0 {
		seed = time.Now().UTC().UnixNano()
	}

	// get screen dimensions from flag otherwise from actual terminal size
	// TODO: validate some minimal size
	width := c.Int("width")
	height := c.Int("height")
	if width <= 0 || height <= 0 {
		tw, th, err := terminal.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			return errors.Wrap(err, "Unable to get terminal size. Set the size manually with --width and --height flags.")
		}
		if width <= 0 {
			width = tw
		}
		if height <= 0 {
			height = th
		}
	}

	// create new game
	game := gorched.NewGame(gorched.GameOptions{
		Width:       width,
		Height:      height,
		Seed:        seed,
		PlayerCount: 2,
		Fps:         c.Int("fps"),
		ASCIIOnly:   c.Bool("ascii-only"),
		LowColor:    c.Bool("low-color"),
		BrowserMode: c.Bool("browser"),
		Debug:       c.Bool("debug"),
	})

	// load demo if requested
	demoPath := c.String("demo")
	if demoPath != "" {
		demo, err := demo.LoadFromFile(demoPath)
		if err != nil {
			return fmt.Errorf("Unable to load demo from file '%s': %w", demoPath, err)
		}
		game.Engine().Screen().AddEntity(demo)
	}

	// start game
	game.Start()

	// greet player at the end
	fmt.Println("Thank you for playing GOrched !")
	fmt.Printf("Your initial seed was: %d\n", game.InitialSeed())
	fmt.Printf("Your last round seed was: %d\n", game.LastSeed())

	// successful finish
	return nil
}
