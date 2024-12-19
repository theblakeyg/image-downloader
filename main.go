package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// create a waitgroup that keep the program alive until the goroutines have finished
var wg sync.WaitGroup

var csv_location = "" //assumption of 1 column of the URLs

var destination_path = "TEST"

type errorLog struct {
	errorString string
	fileName    string
}

var errors []errorLog
var errorsMutex sync.Mutex

func main() {
	//let's open a connection to the file
	f, err := os.Open(csv_location)
	if err != nil {
		log.Fatal(err)
	}

	//remember to close the file when we're done
	defer f.Close()

	//create channel for the urls to be processed from
	urls := make(chan string)

	for w := 1; w <= 50; w++ {
		wg.Add(1) //Add 1 to the Waitgroup, indicating we have another goroutine running
		go func() {
			worker(urls)
		}()
	}

	//now let's read the data with a csvReader
	csvReader := csv.NewReader(f)
	urlCount := 0

	//and load all the urls into the channel
	go func() {
		defer close(urls)

		// Read and discard the header row
		if _, err := csvReader.Read(); err != nil {
			log.Fatal("Error reading header row: ", err)
		}

		for {
			rec, err := csvReader.Read()
			//if we reach the end of file (EOF) then leave the for loop
			if err == io.EOF {
				fmt.Println("End of file")
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			//rec[0] bceause it reads each row as potential multiple columns into a slice, and we just want the first and only one
			urls <- rec[0]
			urlCount++
			fmt.Printf("\rURLs found: %d", urlCount)

		}

		fmt.Printf("\nTotal URLs sent: %d\n", urlCount) // After loop finishes, print the total count
	}()

	wg.Wait() //wait until all goroutines have finished

	printErrors()
}

func printErrors() {
	for _, log := range errors {
		fmt.Println("Error: ", log.errorString, "FileName: ", log.fileName)
	}
}

func worker(urls <-chan string) {
	defer wg.Done() //This signifies to the Waitgroup that 1 of the goroutines it's keeping track of has finished
	for url := range urls {
		//Here is where we want to download the picture
		downloadPicture(url)
	}
}

func downloadPicture(url string) {
	//get the file name
	fileName := getFileName(url)

	newFilePath := filepath.Join(destination_path, fileName)

	//create the file to write to
	img, err := os.Create(newFilePath)
	if err != nil {
		logError(err.Error(), fileName)
		return
	}
	defer img.Close()

	//build the http request
	resp, err := http.Get(url)
	if err != nil {
		logError(err.Error(), fileName)
		return
	}
	defer resp.Body.Close()

	// Check the status code to ensure the image was fetched correctly
	if resp.StatusCode != http.StatusOK {
		logError(fmt.Sprintf("Failed to fetch image, status code: %d", resp.StatusCode), fileName)
		return // Continue processing other URLs even if this one fails
	}

	//save the response (image) to the file location
	b, err := io.Copy(img, resp.Body)
	if err != nil {
		logError(err.Error(), fileName)
		return
	}
	fmt.Println("File size: ", b)

}

func getFileName(url string) string {
	splitUrl := strings.Split(url, "/")
	return splitUrl[len(splitUrl)-1]
}

func logError(errorString, fileName string) {
	errorsMutex.Lock()
	defer errorsMutex.Unlock()

	newError := errorLog{
		errorString: errorString,
		fileName:    fileName,
	}
	errors = append(errors, newError)
}
