// this version is for trying to use flag package to parse the arguments
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

// this struct is for record what user want to do
type config struct {
	numTimes int
	name     string
}

// this function is for recieving user output, bytes.Buffer and os.Stdout can as parameter
// in this function, if you want to check user's input separate, you can create the bytes.Buffer
func getName(r io.Reader, w io.Writer) (string, error) {
	msg := "please input your name, press enter to process:\n"
	fmt.Fprintf(w, "%s", msg)

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	name := scanner.Text()
	if len(name) == 0 {
		return "", errors.New("blank value, this feild should input your name")
	}
	return name, nil

}

var errInvaildPosArgSpecified = errors.New("More than one positional argument specified")

// parse the Arguments what user input.
func parseAgrs(w io.Writer, args []string) (config, error) {
	c := config{}
	fs := flag.NewFlagSet("greeter", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.Usage = func() {
		var usageString = `
	A greeter application which prints the name you entered  a specified
	number of times.

	Usage of %s:`
		fmt.Fprintf(w, usageString, fs.Name())
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}
	fs.IntVar(&c.numTimes, "n", 0, "Number of times to greet")
	err := fs.Parse(args)
	if err != nil {
		return c, err
	}
	//if the quantity of the argument behind the flag
	if fs.NArg() > 1 {
		return c, errInvaildPosArgSpecified
	}
	if fs.NArg() == 1 {
		c.name = fs.Arg(0)
	}
	return c, nil
}

// vaildate user's input
func vaildateArgs(c config) error {
	if !(c.numTimes > 0) {
		return errors.New("invaild number of output times")
	}
	return nil
}

// this is the entrance of executing the scripts' actions
func runCmd(r io.Reader, w io.Writer, c config) error {
	var err error
	if len(c.name) == 0 {
		c.name, err = getName(r, w)
		if err != nil {
			return err
		}
	}
	greetUser(c, w)
	return nil
}

// implement function
func greetUser(c config, w io.Writer) {
	msg := fmt.Sprintf("Nice to meet you, %s\n", c.name)
	for i := 0; i < c.numTimes; i++ {
		fmt.Fprint(w, msg)
	}

}

func main() {
	// 1. try parse the script arguments
	c, err := parseAgrs(os.Stderr, os.Args[1:])
	if err != nil {
		if errors.Is(err, errInvaildPosArgSpecified) {
			fmt.Fprintln(os.Stdout, err)
		}
		os.Exit(1)
	}

	// 2. check the argument if it's illegal
	err = vaildateArgs(c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}

	// 3. execute the implement
	err = runCmd(os.Stdin, os.Stdout, c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
