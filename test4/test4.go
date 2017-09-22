/*
Test #3
Modify the codebase to run the following filters:
1 Filter out cheap properties (anything under 400k)
2 Filter out properties that are avenues, crescents, or places (AVE, CRES, PL) cos
those guys are just pretentious...
3 Filter out every 10th property (to keep our users on their toes!)
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
  "strconv"
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

func StrToInt(str string) (int, error) {
    nonFractionalPart := strings.Split(str, ".")
    return strconv.Atoi(nonFractionalPart[0])
}

func filter(ch1 chan []Property, divided []Property)  {
	var filterRes []Property
	for count , prop := range divided {
		ok, _ := regexp.MatchString("(?i)AVE|CRES|PL", prop.Address)
		value, _ := StrToInt(prop.Value)
		if !ok && value > 400000 && count+1 != 10 {
			filterRes = append(filterRes, prop)
		}
	}
	ch1 <- filterRes 
}




func main() {
	lines, err := readLines("../properties.txt")
	if err != nil {
		fmt.Println("readLines: %s", err)
	}
	
	// iterate through lines and create an array with a list of properties without duplicates 
	var property []Property
	//Modify the code in case of duplicates to use the first encountered record.
	for _, ln := range lines {
		// Splits the string on tab spaces and concatenates it with the value to create an array
		r :=  regexp.MustCompile(`(.*?)\t|\d*`)
		line := r.FindAllString(ln, -1)
		//validate addresses
		if len(line) >= 5 {
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
	
	//chunks
	chunkSize := 10
	var filterRes []Property 
	
	
	for i := 0; i < len(property); i += chunkSize {
		end := i + chunkSize

		if end > len(property) {
			end = len(property)
		}
		divided := property[i:end]
		
		ch := make(chan []Property)
		
		go filter(ch, divided)
		
		
		for _, f := range <-ch {
			filterRes = append(filterRes,f)
		
		}
	}
	
	filterResJson , _	:= json.Marshal(filterRes)
   	fmt.Println("Filtered Result \n", jsonPrettyPrint(string(filterResJson)))
}
