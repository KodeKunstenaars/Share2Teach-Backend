package convertToPDF

import (
	"bufio"         // allows text files to be read line by line
	"fmt"           // used to print to console
	"os"            // used to interact with the operating system (allows the code to open files)
	"path/filepath" // allows for the manipulation of file paths

	"github.com/jung-kurt/gofpdf/v2" // Used for the conversion to pdf
)

func convertTextToPDF(inputFile string) error {

	//New PDF doc is created
	pdf := gofpdf.New("P", "mm", "A4", "")
	// blank page is added that will contain the content of the textfile
	pdf.AddPage()

	// PDF font is set
	pdf.SetFont("Arial", "", 12)

	// The text file is opened
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("Failed to open input file: %w", err)
	}

	// ensures that the file is closed, even in the event of an error
	defer file.Close()

	// the data of the file is read into a buffer
	scanner := bufio.NewScanner(file)

	// y-coordinate for the text is initialised
	y := 10.0

	// file is read line by line and added to the pdf document
	for scanner.Scan() {
		line := scanner.Text()
		pdf.Text(10, y, line)
		y += 10
	}

	// checks for errors during scanning
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading the file: %w", err)
	}

	// the pdf file name is generated with the .pdf extention
	outputfile := filepath.Base(inputFile)
	outputfile = outputfile[:len(outputfile)-len(filepath.Ext(outputfile))] + ".pdf"

	// The pdf is saved to a file
	err = pdf.OutputFileAndClose(outputfile)
	if err != nil {
		return fmt.Errorf("Error saving the file: %w", err)
	}

	return nil
}

/*
func main() {
	err := convertTextToPDF("Test.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("PDF created successfully!")
}
*/
