package utils

import (
	"bufio"
	"log"
	"os"
)

func GetUserInput(message []byte) (string, error) {

	if message[len(message)-1] == '\n' {
		os.Stdin.Write(message)
	} else {
		os.Stdin.Write(message)
		os.Stdin.Write([]byte("\n"))
	}

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	var userInput string = line[:len(line)-1]

	return userInput, nil
}
