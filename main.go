package main

import (
	"bytes"
	"html/template"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

// Struktur untuk menampung data JSON
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
	r := gin.Default()

	r.POST("/generate", func(c *gin.Context) {
		var resume Resume
		if err := c.BindJSON(&resume); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Load HTML template
		tmpl, err := template.ParseFiles("template.html")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load template"})
			return
		}

		// Render HTML template with data
		var htmlBuffer bytes.Buffer
		if err := tmpl.Execute(&htmlBuffer, resume); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to render template"})
			return
		}

		// Create PDF using wkhtmltopdf
		cmd := exec.Command("wkhtmltopdf", "-", "-")
		cmd.Stdin = &htmlBuffer
		pdfBuffer, err := cmd.Output()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate PDF"})
			return
		}

		// Send PDF as response
		c.Data(http.StatusOK, "application/pdf", pdfBuffer)
	})

	r.Run(":8001")
}
