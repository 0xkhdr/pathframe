package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/0xkhdr/pathframe/internal/app"
)

// Run executes the Pathframe command and returns its process exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("pathframe", flag.ContinueOnError)
	flags.SetOutput(stderr)
	help := flags.Bool("help", false, "show help")
	version := flags.Bool("version", false, "show version")
	flags.Usage = func() { usage(stderr) }

	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *help {
		usage(stdout)
		return 0
	}
	if *version {
		fmt.Fprintln(stdout, app.Version)
		return 0
	}
	if flags.NArg() > 0 {
		fmt.Fprintf(stderr, "pathframe: unknown command %q\n", flags.Arg(0))
		return 2
	}

	fmt.Fprintln(stdout, "Pathframe")
	fmt.Fprintln(stdout, "")
	fmt.Fprintln(stdout, "Status: not initialized")
	fmt.Fprintln(stdout, "Next: workflow commands are not available in Stage 0.")
	return 0
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "Usage: pathframe [--help] [--version]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Options:")
	fmt.Fprintln(w, "  --help     show help")
	fmt.Fprintln(w, "  --version  show version")
}
