package pdfconvert

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf/v2"
)

// reads the file into memory and converts it to a PDF in-memory
func ConvertTextToPDF(filebytes []byte, filename string) ([]byte, error) {

	//New PDF doc is created
	pdf := gofpdf.New("P", "mm", "A4", "")
	// blank page is added that will contain the content of the textfile
	pdf.AddPage()

	// PDF font is set
	pdf.SetFont("Arial", "", 12)

	// a buffer is used to read the file line by line
	buffer := bytes.NewBuffer(filebytes)
	scanner := bufio.NewScanner(buffer)

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
		return nil, fmt.Errorf("Error Reading the file: %w", err)
	}

	// pdf is put into a byte buffer
	var pdfBuffer bytes.Buffer
	if err := pdf.Output(&pdfBuffer); err != nil {
		return nil, fmt.Errorf("Error generating the PDF: %w", err)
	}

	return pdfBuffer.Bytes(), nil
}
