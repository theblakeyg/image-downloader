# Image Downloader
A program to download images from a list of urls

## How to
1. Create a csv with one column and a header the holds all of the URLs
2. In main.go set csv_location to the relative file path of your csv
3. Set the destination_path to the relative folder path of where you want the images downloaded to
4. Either compile the executable and run or "go run main.go" to start the process

## Technical Details
As we read the CSV line by line we put the URL onto the channel with a default of 50 workers as receivers.<br/>
The next available worker takes the newest URL off the channel and downloads the image from the URL and saves it to the destination_path.<br/>
Errors are logged at possible points of failure without stopping the program and printed at the end of the run.<br/>

## ToDo
[] Defaults for csv_location, destination_path, workers<br/>
[] Capture total number of URLs before we start adding to the channel for UX<br/>
[] Display real-time progress bar in command line for UX<br/>
[] Display real-time number of erros in command line for UX<br/>
[] Display real-time number of images processed in command line for UX<br/>
[] Disaply real-time data transfer number in command line for UX<br/>
[] Log errors to a file for persistence<br/>
[] Be able to call from CLI with parameters for --csv-location, --destination-path, --workers<br/>
