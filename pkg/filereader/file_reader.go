package filereader

import (
	"bufio"
	"os"
)

func ReadLines(filename string) ([]string, error) {
	readFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	lines := make([]string, 0)

	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}

	err = readFile.Close()
	return lines, err
}
