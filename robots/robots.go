package robots

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func FindRobotsTxtURLs(inputFile string, outputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	defer file.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer output.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		if strings.Contains(url, "robots.txt") {
			output.WriteString(url + "\n") 
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}
	return nil
}

