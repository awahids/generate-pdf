package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type Resume struct {
	Basics struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Location struct {
			Address string `json:"address"`
		} `json:"location"`
		Website string `json:"website"`
		Summary string `json:"summary"`
	} `json:"basics"`
	Education []struct {
		Institution string `json:"institution"`
		Area        string `json:"area"`
		StartDate   string `json:"startDate"`
		EndDate     string `json:"endDate"`
	} `json:"education"`
	Work []struct {
		Position   string   `json:"position"`
		Company    string   `json:"company"`
		StartDate  string   `json:"startDate"`
		EndDate    string   `json:"endDate"`
		Highlights []string `json:"highlights"`
	} `json:"work"`
}

func main() {
	// Membaca file JSON
	file, err := os.ReadFile("resume_eng.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// Unmarshal JSON ke dalam struktur Resume
	var resume Resume
	err = json.Unmarshal(file, &resume)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Render template HTML
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		fmt.Println("Error parsing template file:", err)
		return
	}

	// Buat file sementara untuk menyimpan hasil render HTML
	tmpFile, err := os.CreateTemp("", "resume_*.html")
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	err = tmpl.Execute(tmpFile, resume)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	tmpFile.Close()

	// Konversi HTML ke PDF menggunakan wkhtmltopdf
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		fmt.Println("Error creating PDF generator:", err)
		return
	}

	pdfg.AddPage(wkhtmltopdf.NewPage(tmpFile.Name()))

	err = pdfg.Create()
	if err != nil {
		fmt.Println("Error creating PDF:", err)
		return
	}

	err = pdfg.WriteFile("resume.pdf")
	if err != nil {
		fmt.Println("Error writing PDF:", err)
		return
	}

	fmt.Println("PDF successfully created")
}
