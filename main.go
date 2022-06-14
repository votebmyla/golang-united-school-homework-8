package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}

}

func Perform(args Arguments, writer io.Writer) error {

	switch args["operation"] {
	case "add":
		fmt.Println("this is add operation")
	case "list":
		if len(args["fileName"]) == 0 {
			return errors.New("-fileName flag has to be specified")
		}
		listItems(args, writer)
	case "findById":
		findById(args)
	case "remove":
		fmt.Println("this is remove operation")
	default:
		if len(args["operation"]) > 0 {
			return fmt.Errorf("Operation %s not allowed!", args["operation"])
		}
		return errors.New("-operation flag has to be specified")
	}

	return nil
}

type Arguments map[string]string

func parseArgs() Arguments {
	userOperation := flag.String("operation", "", "invalid operation")
	userId := flag.String("id", "-1", "invalid id")
	userItem := flag.String("item", "", "invalid item")
	userFilename := flag.String("fileName", "", "invalid filename")

	flag.Parse()
	var userArgs Arguments = Arguments{
		"operation": *userOperation,
		"id":        *userId,
		"item":      *userItem,
		"fileName":  *userFilename,
	}

	return userArgs
}

func listItems(userArgs Arguments, writer io.Writer) error {
	// fmt.Println(userArgs)
	file, err := os.OpenFile(userArgs["fileName"], os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	buf := make([]byte, 16)
	for {
		length, err := file.Read(buf)
		if err == io.EOF {
			return err
		}
		if err != nil {
			return err
		}
		writer.Write(buf[:length])
	}
}

func findById(args Arguments) error {

	// file, err := os.OpenFile("")
	// for i, v := range args {
	// 	fmt.Println(args[i], v)
	// }

	return nil
}

// func checkOperationFlag(args Arguments) error {
// 	switch args["operation"]{
// 		case ""
// 	}
// }
