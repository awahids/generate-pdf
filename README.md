# Resume PDF Generator

Proyek ini adalah aplikasi untuk menghasilkan resume dalam bentuk PDF menggunakan data JSON. Aplikasi ini menggunakan Gin Gonic sebagai framework web dan wkhtmltopdf untuk membuat PDF dari HTML.

## Instalasi

1. **Clone Repository**
   ```sh
   git clone https://github.com/awahids/generate-pdf.git
   cd generate-pdf
   ```

2. **Install Dependencies**
   - **Gin Gonic**:
     ```sh
     go get -u github.com/gin-gonic/gin
     ```
   - **wkhtmltopdf**:
     ```sh
     brew install wkhtmltopdf
     ```

3. **Download Font Lexend Deca**
   - Pastikan Anda terhubung ke internet untuk mengunduh font dari Google Fonts.

## Menghasilkan PDF tanpa JSON

Anda dapat membuat PDF langsung dari template HTML dengan data statis.

1. **Buat File HTML Template**:
   Buat file `template.html` dengan konten sebagai berikut:

   ```html
   <!DOCTYPE html>
   <html lang="id">
   <head>
       <meta charset="UTF-8">
       <meta name="viewport" content="width=device-width, initial-scale=1.0">
       <title>Resume</title>
       <link href="https://fonts.googleapis.com/css2?family=Lexend+Deca&display=swap" rel="stylesheet">
       <style>
           body {
               font-family: 'Lexend Deca', sans-serif;
               margin: 40px;
               color: #333;
           }
           .resume {
               width: 800px;
               margin: auto;
               padding: 20px;
           }
           .header {
               text-align: center;
               margin-bottom: 20px;
           }
           .header h1 {
               margin: 0;
               font-size: 36px;
           }
           .header p {
               margin: 5px 0;
               font-size: 14px;
           }
           .section {
               margin-bottom: 20px;
           }
           .section-title {
               font-size: 24px;
               color: #333;
               border-bottom: 2px solid #666;
               margin-bottom: 10px;
           }
           .item {
               margin-bottom: 10px;
           }
           .item-label {
               font-weight: bold;
           }
           .item-detail {
               margin-left: 20px;
           }
       </style>
   </head>
   <body>
       <div class="resume">
           <div class="header">
               <h1>John Doe</h1>
               <p>johndoe@example.com</p>
               <p>+628123456789</p>
               <p>Jalan Raya No. 123, Jakarta</p>
               <p>www.johndoe.com</p>
           </div>
           <div class="section">
               <div class="section-title">Summary</div>
               <p>Experienced backend developer with expertise in creating scalable web applications.</p>
           </div>
           <div class="section">
               <div class="section-title">Pendidikan</div>
               <div class="item">
                   <span class="item-label">Institusi:</span> Universitas X
                   <div class="item-detail"><span class="item-label">Bidang Studi:</span> Teknik Informatika</div>
                   <div class="item-detail"><span class="item-label">Tanggal Mulai:</span> 2015-09-01</div>
                   <div class="item-detail"><span class="item-label">Tanggal Selesai:</span> 2019-06-30</div>
               </div>
           </div>
           <div class="section">
               <div class="section-title">Pengalaman Kerja</div>
               <div class="item">
                   <span class="item-label">Posisi:</span> Backend Developer
                   <div class="item-detail"><span class="item-label">Perusahaan:</span> Your Company</div>
                   <div class="item-detail"><span class="item-label">Tanggal Mulai:</span> 2019-07-01</div>
                   <div class="item-detail"><span class="item-label">Tanggal Selesai:</span> Present</div>
                   <div class="item-detail"><span class="item-label">Highlight:</span> Developed REST APIs using NestJS</div>
                   <div class="item-detail"><span class="item-label">Highlight:</span> Managed database migrations and optimizations</div>
               </div>
           </div>
       </div>
   </body>
   </html>
   ```

## Menggunakan REST API untuk Generate PDF

REST API ini dibuat di branch `rest-api`. Anda dapat membuat PDF dari data JSON menggunakan endpoint `/generate`.

1. **Checkout ke Branch `rest-api`**:
   ```sh
   git checkout rest-api
   ```

2. **Buat File `main.go`**:
   Buat file `main.go` dengan konten sebagai berikut:

   ```go
   package main

   import (
       "bytes"
       "html/template"
       "io/ioutil"
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

       r.Run(":8080")
   }
   ```

3. **Jalankan API**:
   Jalankan server dengan perintah berikut:
   ```sh
   go run main.go
   ```

4. **Uji API**:
   Gunakan alat seperti Postman atau curl untuk mengirim permintaan POST dengan JSON body:

   ```json
   {
       "basics": {
           "name": "John Doe",
           "email": "johndoe@example.com",
           "phone": "+628123456789",
           "location": {
               "address": "Jalan Raya No. 123, Jakarta"
           },
           "website": "www.johndoe.com",
           "summary": "Experienced backend developer with expertise in creating scalable web applications."
       },
       "education": [
           {
               "institution": "Universitas X",
               "area": "Teknik Informatika",
               "startDate": "2015-09-01",
               "endDate": "2019-06-30"
           }
       ],
       "work": [
           {
               "position": "Backend Developer",
               "company": "Your Company",
               "startDate": "2019-07-01",
               "endDate": "Present",
               "highlights": [
                   "Developed REST APIs using NestJS",
                   "Managed database migrations and optimizations"
               ]
           }
       ]
   }
   ```

   Contoh perintah curl:
   ```sh
   curl -X POST http://localhost:8080/generate -H "Content-Type: application/json" -d @resume.json --output resume.pdf
   ```
