package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}

}

func Perform(args Arguments, writer io.Writer) error {
	// fmt.Println(args)
	if len(args["fileName"]) == 0 {
		return errFilename
	}
	switch args["operation"] {
	case "add":
		if len(args["item"]) == 0 {
			return errItem
		}
		addItem(args, writer)
	case "list":
		listItems(args, writer)
	case "findById":
		if len(args["id"]) == 0 {
			return errId
		}
		findById(args, writer)
	case "remove":
		if len(args["id"]) == 0 {
			return errId
		}
		removeItem(args, writer)
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
	userId := flag.String("id", "", "invalid id")
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
	file, err := os.OpenFile(userArgs["fileName"], os.O_RDONLY, 0444)
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

func addItem(args Arguments, writer io.Writer) error {
	jsonFile, err := os.Open(args["fileName"])
	if err != nil {
		newJsonFile, err := os.Create(args["fileName"])
		if err != nil {
			return err
		}
		defer newJsonFile.Close()
		var userDataInput User
		err = json.Unmarshal([]byte(args["item"]), &userDataInput)
		if err != nil {
			return err
		}
		var userDataInputArr []User
		userDataInputArr = append(userDataInputArr, userDataInput)
		userDataInputArrByte, err := json.Marshal(userDataInputArr)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(args["fileName"], userDataInputArrByte, 0644)
		if err != nil {
			return err
		}
	}
	defer jsonFile.Close()
	jsonFileByte, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	var jsonFileArr []User
	err = json.Unmarshal(jsonFileByte, &jsonFileArr)
	if err != nil {
		return err
	}
	//
	for _, v := range jsonFileArr {
		var itemData User
		err := json.Unmarshal([]byte(args["item"]), &itemData)
		if err != nil {
			return err
		}
		if itemData.Id == v.Id {
			s := fmt.Sprintf("Item with id %s already exists", itemData.Id)
			writer.Write([]byte(s))
		}
	}
	//

	return nil
}

func removeItem(args Arguments, writer io.Writer) error {
	jsonFile, err := os.Open(args["fileName"])
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	jsonFileByte, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	var jsonFileDataArr []User
	err = json.Unmarshal(jsonFileByte, &jsonFileDataArr)
	if err != nil {
		return err
	}

	for i, v := range jsonFileDataArr {
		if v.Id == args["id"] {
			jsonFileDataArr = append(jsonFileDataArr[:i], jsonFileDataArr[i+1:]...)
			byteValue, err := json.Marshal(jsonFileDataArr)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(args["fileName"], byteValue, 0644)
			if err != nil {
				return err
			}
		}
	}

	s := fmt.Sprintf("Item with id %s not found", args["id"])
	writer.Write([]byte(s))
	return nil
}

func findById(args Arguments, writer io.Writer) error {
	jsonFile, err := os.OpenFile(args["fileName"], os.O_RDONLY, 0444)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	var users Users
	json.Unmarshal(byteValue, &users)
	for _, v := range users {
		if v.Id == args["id"] {
			data, err := json.Marshal(v)
			if err != nil {
				return err
			}
			writer.Write(data)
		}
	}
	return nil
}

var errFilename error = errors.New("-fileName flag has to be specified")
var errId error = errors.New("-id flag has to be specified")
var errItem error = errors.New("-item flag has to be specified")

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type Users []User
