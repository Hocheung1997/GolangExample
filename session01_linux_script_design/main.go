package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

var argumentCount int = 1

// define the discription of the scripts
var usageString = fmt.Sprintf(`Usage :%s <integer> [-h|--help]

A greeter application which prints the name you entered <integer>
number of times.

`, os.Args[0])

func promptUsage(w io.Writer) {
	fmt.Fprint(w, usageString)
}

// this struct is for record what user want to do
type config struct {
	numTimes   int
	printUsage bool
}

// this function is for recieving user output, bytes.Buffer and os.Stdout can as parameter
// in this function, if you want to check user's input separate, you can create the bytes.Buffer
func getName(r io.Reader, w io.Writer) (string, error) {
	msg := "pleae input your name, press enter to process:\n"
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

// parse the Arguments what user input.
func parseAgrs(args []string) (config, error) {
	var numTimes int
	var err error
	c := config{}
	if len(args) != argumentCount {
		return config{}, errors.New("invaild number of arguments")
	}

	if args[0] == "-h" || args[0] == "--help" {
		c.printUsage = true
		return c, nil
	}
	numTimes, err = strconv.Atoi(args[0])
	if err != nil {
		return c, err
	}
	c.numTimes = numTimes

	return c, nil

}

// vaildate user's input
func vaildateArgs(c config) error {
	if c.printUsage {
		return nil
	}
	if !(c.numTimes > 0) {
		return errors.New("invaild number of output times")
	}
	return nil
}

// this is the entrance of executing the scripts' actions
func runCmd(r io.Reader, w io.Writer, c config) error {
	if c.printUsage {
		fmt.Fprint(w, "usage example: ./application 6")
		return nil
	}
	name, err := getName(r, w)
	if err != nil {
		return err
	}
	greetUser(c, name, w)
	return nil
}

// implement function
func greetUser(c config, name string, w io.Writer) {
	msg := fmt.Sprintf("Nice to meet you, %s\n", name)
	for i := 0; i < c.numTimes; i++ {
		fmt.Fprint(w, msg)
	}

}

func main() {
	// 1. try parse the script arguments
	c, err := parseAgrs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		promptUsage(os.Stdout)
		os.Exit(1)
	}

	// 2. check the argument if it's illegal
	err = vaildateArgs(c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		promptUsage(os.Stdout)
		os.Exit(1)
	}

	// 3. execute the implement
	err = runCmd(os.Stdin, os.Stdout, c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
