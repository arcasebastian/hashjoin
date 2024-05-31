// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Hello is a simple hello, world demonstration web server.
//
// It serves version information on /version and answers
// any other request like /name by saying "Hello, name!".
//
// See golang.org/x/example/outyet for a more sophisticated server.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/bxcodec/faker/v3"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: helloserver [options]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	addr = flag.String("addr", "localhost:8081", "address to serve")
)

func main() {
	// Parse flags.
	flag.Usage = usage
	flag.Parse()

	// Parse and validate arguments (none).
	args := flag.Args()
	if len(args) != 0 {
		usage()
	}

	// Register handlers.
	// All requests not otherwise mapped with go to greet.
	// /version is mapped specifically to version.
	http.HandleFunc("/", greet)
	http.HandleFunc("/version", version)

	log.Printf("serving http://%s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func version(w http.ResponseWriter, r *http.Request) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		http.Error(w, "no build information available", 500)
		return
	}

	fmt.Fprintf(w, "<!DOCTYPE html>\n<pre>\n")
	fmt.Fprintf(w, "%s\n", html.EscapeString(info.String()))
}

func greet(w http.ResponseWriter, r *http.Request) {
	joinType := strings.Trim(r.URL.Path, "/")
	if joinType == "" {
		joinType = "inner"
	}

	resultSet := product(leftData(), rightData(), "assetId", joinType)
	jsonData, err := json.Marshal(resultSet)
	if err != nil {
		log.Print("Error")
		return
	}
	fmt.Fprintf(w, "<!DOCTYPE html>\n<pre>\n")
	fmt.Fprintf(w, "SELECT * FROM asset %s JOIN report ON asset.assetId = report.assetId\n", strings.ToUpper(joinType))
	fmt.Fprintf(w, "%s\n", html.EscapeString(string(jsonData)))
}

func product(leftList []map[string]any, rightList []map[string]any, joinableKey string, joinType string) []map[string]any {
	switch joinType {
	case "inner":
		return innerJoin(leftList, rightList, joinableKey)
	case "right":
		return outerJoin(leftList, rightList, joinableKey)
	case "left":
		return outerJoin(rightList, leftList, joinableKey)
	}
	return nil
}

func innerJoin(buildList []map[string]any, probeList []map[string]any, joinableKey string) []map[string]any {
	buildTableMap := buildTable(buildList, joinableKey)
	var result []map[string]any
	for _, item := range probeList {
		key := item[joinableKey]
		find := buildTableMap[key]
		if find != nil {
			for i := 0; i < len(find); i++ {
				result = append(result, joinRow(item, find[i]))
			}
		}
	}
	return result
}

func outerJoin(buildList []map[string]any, probeList []map[string]any, joinableKey string) []map[string]any {
	emptyRow := makeEmptyRow(buildList[0])
	buildTableCol := buildTable(buildList, joinableKey)
	var result []map[string]any
	for _, item := range probeList {
		key := item[joinableKey]
		find := buildTableCol[key]
		if find != nil {
			for i := 0; i < len(find); i++ {
				result = append(result, joinRow(item, find[i]))
			}
		} else {
			result = append(result, joinRow(emptyRow, item))
		}
	}
	return result
}

/*func fullJoin(leftData []map[string]any, rightData []map[string]any, joinableKey string) []map[string]any {
	emptyRowLeft := makeEmptyRow(leftData[0])
	emptyRowRight := makeEmptyRow(rightData[0])

	var result []map[string]any
	var leftPointer, rightPointer = 0, 0
	for leftPointer < len(leftData) && rightPointer < len(rightData) {
		leftItem := leftData[leftPointer]
		rightItem := rightData[rightPointer]

		if leftItem[joinableKey] == rightItem[joinableKey] {
			result = append(result, joinRow(leftItem, rightItem))
			leftPointer++
			rightPointer++
		} else if leftItem[joinableKey] < rightItem[joinableKey] {
			result = append(result, joinRow(emptyRowLeft, rightItem))
			rightPointer++
		} else {
			result = append(result, joinRow(leftItem, emptyRowRight))
			leftPointer++
		}

	}
	return result
}*/

func buildTable(rightList []map[string]any, joinableKey string) map[any][]map[string]any {
	buildTable := make(map[any][]map[string]any)

	for _, item := range rightList {
		key := item[joinableKey]
		buildTable[key] = append(buildTable[key], item)
	}
	return buildTable
}

func makeEmptyRow(row map[string]any) map[string]any {
	emptyRow := make(map[string]any)
	for key := range row {
		emptyRow[key] = nil
	}
	return emptyRow
}

/*func sortDataset(list []map[string]any, key string) []map[string]any {
	keysToSort := make([]any, len(list))
	keydTable := buildTable(list, key)
	for i := 0; i < len(list); i++ {
		keysToSort = append(keysToSort, list[i][key])
	}
	sort.Strings(keysToSort)
	sortedList := make([]map[string]any, len(list))
	for i := 0; i < keysToSort; i++ {
		sortedList = append(sortedList, keydTable[[i]])
	}
	return sortedList

}*/

func joinRow(leftRow map[string]any, rightRow map[string]any) map[string]any {
	joinedRow := make(map[string]any)
	for key, value := range leftRow {
		joinedRow[key] = value
	}
	for key, value := range rightRow {
		joinedRow[key] = value
	}
	return joinedRow
}

func rightData() []map[string]any {
	var list []map[string]any
	faker.Name()

	for i := 0; i < 5000; i++ {
		element := make(map[string]any)
		element["assetId"] = rand.Intn(2000)
		element["latitude"] = faker.Latitude()
		element["longitude"] = faker.Longitude()
		list = append(list, element)
	}

	return list
}

func leftData() []map[string]any {
	var list []map[string]any

	for i := 0; i < 1500; i++ {
		element := make(map[string]any)
		element["assetId"] = i
		element["driverName"] = faker.Name()
		element["domain"] = faker.DomainName()
		list = append(list, element)
	}
	return list
}
