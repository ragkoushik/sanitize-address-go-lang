/*
Test #3
Instead of inserting the last record, make sure that no duplicates are entered at all.
i.e., if there are duplicate records, do not insert any records.
Author: Raghu Koushik 			
*/ 

package main

import (
  "bufio"
  "fmt"
  "encoding/json"
  "os"
  "regexp"
  "bytes"
  "strings"
  "unicode"
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

func SpaceMap(str string) string {
    return strings.Map(func(r rune) rune {
        if unicode.IsSpace(r) {
            return -1
        }
        return r
    }, str)
}

func in_array(val string, array []Property) (skip bool, index int) {
    skip = false
    index = -1;

    for i, v := range array {
        if SpaceMap(val) == SpaceMap(v.Address + v.Town + v.Date) {
            index = i
            skip = true
            return
        }   
    }

    return
}

func main() {
	lines, err := readLines("../properties.txt")
	if err != nil {
		fmt.Println("readLines: %s", err)
	}
	// iterate through lines and create an array with a list of properties without duplicates 
	var property []Property
	//Modify the code in case of duplicates to use the first encountered record.
	for i, ln := range lines {
		// Splits the string on tab spaces and concatenates it with the value to create an array
		r :=  regexp.MustCompile(`(.*?)\t|\d*`)
		line := r.FindAllString(ln, -1)
		if i > 0 && len(line) >= 5 {
			//record the property after inserting 
			// address, town and date needs to be unique - creating a unique key by concatenating Address + Town + Date strings
			
			skip,index := in_array( line[1]+line[2]+line[3] , property )
						
			if !skip {
				property = append(property, Property{
					Id: 		line[0],
					Address: 	line[1],
					Town: 		line[2],
					Date:	 	line[3],
					Value:		line[4],
				})
			} else {
				property = append(property[:index], property[index+1:]...)
			}
		} 
	}
	//Print Result
	propertyJson, _ := json.Marshal(property)
	fmt.Println(jsonPrettyPrint(string(propertyJson)))
}
