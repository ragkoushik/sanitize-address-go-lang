/*
Test #1
Write a routine that implements the logic described above. Use the attached CSV as input.

In the case of duplicates, use the last encountered record.
NOTE: A duplicate is a row that has the same address and same date. The ID is irrelevant.

Print the list.
Author: Raghu Koushik 
*/ 

package main

import (
  "bufio"
  "fmt"
  "encoding/json"
  "os"
  "regexp"
  "strings"
  "unicode"
  "bytes"
)

type Property struct {
    Id 			string 
    Address     string
	Town		string
	Date	 	string
	Value		string
}

// readLines reads a whole file into memory
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
  return lines, scanner.Err()
}

// Create a Json pretty print
func jsonPrettyPrint(in string) string {
    var out bytes.Buffer
    err := json.Indent(&out, []byte(in), "", "\t")
    if err != nil {
        return in
    }
    return out.String()
}

//removes white spaces
func SpaceMap(str string) string {
    return strings.Map(func(r rune) rune {
        if unicode.IsSpace(r) {
            return -1
        }
        return r
    }, str)
}


func main() {
	lines, err := readLines("../properties.txt")
	if err != nil {
		fmt.Println("readLines: %s", err)
	}
	// iterate through lines and create an array with a list of properties without duplicates 
	var property []Property
	
	//filter duplicates
	PropExists := make(map[string]bool)
	//In the case of duplicates, use the last encountered record. - iterating in the reverse order to pick the last record first
	for i := len(lines)-1; i >= 0; i-- {
		// Splits the string on tab spaces and concatenates it with the value to create an array
		r :=  regexp.MustCompile(`(.*?)\t|\d*`)
		line := r.FindAllString(lines[i], -1)
		if i > 0 && len(line) >= 5 && !PropExists[SpaceMap(line[1]+line[2]+line[3]) ] {
			//record the property after inserting 
			// address, town and date needs to be unique - creating a unique key by concatenating Address + Town + Date strings
			PropExists[SpaceMap(line[1]+line[2]+line[3])] = true
			
			property = append(property, Property{
					Id: 		line[0],
					Address: 	line[1],
					Town: 		line[2],
					Date:	 	line[3],
					Value:		line[4],
				})
		} 
	}
	//Print Result
	propertyJson, _ := json.Marshal(property)
	fmt.Println(jsonPrettyPrint( string(propertyJson)))
}
