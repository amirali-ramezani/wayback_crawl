package javascript

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)


func FindJSURLs(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("error opening file: %w", err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        url := scanner.Text()
        if strings.Contains(url, ".js") {
            fmt.Println(url) 
        }
    }

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("error reading file: %w", err)
    }
    return nil
}
