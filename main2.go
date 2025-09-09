package main

// import (
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"path/filepath"

// 	"github.com/pdfcpu/pdfcpu/pkg/api"
// 	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
// 	qrcode "github.com/skip2/go-qrcode"
// )

// type SignatureInfo struct {
// 	OrgName string
// 	QRText  string
// }

// type QRStampingConfig struct {
// 	Position string  `yaml:"position"` // "br" (bottom right), "bl" (bottom left), "c" (center), etc.
// 	OffsetX  int     `yaml:"offset_x"` // X offset in points
// 	OffsetY  int     `yaml:"offset_y"` // Y offset in points
// 	Scale    float64 `yaml:"scale"`    // Scale factor (0.1 = 10% of original size)
// 	Rotation int     `yaml:"rotation"` // Rotation in degrees (0 = no rotation)
// 	Page     string  `yaml:"page"`     // Page specification ("l" = last page, "1" = first page, etc.)
// }

// func main() {
// 	inputPDF := "input.pdf"
// 	outputPDF := "output.pdf"

// 	signatures := []SignatureInfo{
// 		{"Company A", "Imzo A"},
// 		{"Company B", "Imzo B"},
// 		{"Company C", "Imzo C"},
// 		{"Company D", "Imzo D"},
// 	}

// 	conf := pdfcpu.NewDefaultConfiguration()

// 	// Temp directory for QR codes
// 	tempDir := "./tempQR"
// 	os.Mkdir(tempDir, 0755)

// 	// Coordinates and page settings
// 	pageWidth := 595.0  // A4 width (points)
// 	pageHeight := 842.0 // A4 height (points)
// 	margin := 50.0
// 	qrSize := 100.0
// 	textHeight := 20.0
// 	itemHeight := qrSize + textHeight + 10.0 // QR + text + padding
// 	startY := pageHeight - margin
// 	currentY := margin + itemHeight
// 	currentPage := 1

// 	for i, sig := range signatures {
// 		// QR file path
// 		qrFile := fmt.Sprintf("%s/qr_%d.png", tempDir, i)
// 		err := qrcode.WriteFile(sig.QRText, qrcode.Medium, int(qrSize), qrFile)
// 		if err != nil {
// 			panic(err)
// 		}

// 		// Agar joy qolmasa → yangi sahifa
// 		if currentY+itemHeight > startY {
// 			currentPage++
// 			err := api.InsertPagesFile(outputPDF, "", []string{fmt.Sprintf("%d", currentPage)}, false, conf)
// 			if err != nil {
// 				panic(err)
// 			}
// 			currentY = margin
// 		}

// 		// QR joylashtirish
// 		pos := fmt.Sprintf("%.2f %.2f", margin, currentY)
// 		err = api.AddImageFile(outputPDF, outputPDF, qrFile, pos, nil, conf)
// 		if err != nil {
// 			panic(err)
// 		}

// 		// Matn joylashtirish (QR ostida)
// 		textPos := fmt.Sprintf("%.2f %.2f", margin, currentY-15)
// 		err = api.AddTextFile(outputPDF, outputPDF, []string{fmt.Sprintf("%d", currentPage)}, textPos, sig.OrgName, conf)
// 		if err != nil {
// 			panic(err)
// 		}

// 		currentY += itemHeight
// 	}

// 	fmt.Println("Bajarildi:", outputPDF)
// }

// func AddQRCodeToPDF(originalPDF []byte, qrCodePNG []byte, config QRStampingConfig) ([]byte, error) {
// 	tempDir := "/tmp"

// 	// Create temporary files
// 	pdfFile := filepath.Join(tempDir, fmt.Sprintf("input_%d.pdf", os.Getpid()))
// 	qrFile := filepath.Join(tempDir, fmt.Sprintf("qr_%d.png", os.Getpid()))
// 	outputFile := filepath.Join(tempDir, fmt.Sprintf("output_%d.pdf", os.Getpid()))

// 	// Clean up temporary files
// 	defer func() {
// 		os.Remove(pdfFile)
// 		os.Remove(qrFile)
// 		os.Remove(outputFile)
// 	}()

// 	// Write input files
// 	if err := os.WriteFile(pdfFile, originalPDF, 0644); err != nil {
// 		return nil, fmt.Errorf("failed to write PDF file: %w", err)
// 	}

// 	if err := os.WriteFile(qrFile, qrCodePNG, 0644); err != nil {
// 		return nil, fmt.Errorf("failed to write QR file: %w", err)
// 	}

// 	// Build pdftk command for stamping
// 	cmd := exec.Command("pdftk", pdfFile, "stamp", qrFile, "output", outputFile)

// 	// Set stamping parameters based on configuration
// 	env := os.Environ()
// 	env = append(env, fmt.Sprintf("PDFTK_STAMP_POSITION=%s", config.Position))
// 	env = append(env, fmt.Sprintf("PDFTK_STAMP_OFFSET_X=%d", config.OffsetX))
// 	env = append(env, fmt.Sprintf("PDFTK_STAMP_OFFSET_Y=%d", config.OffsetY))
// 	env = append(env, fmt.Sprintf("PDFTK_STAMP_SCALE=%.2f", config.Scale))
// 	cmd.Env = env

// 	if output, err := cmd.CombinedOutput(); err != nil {
// 		return nil, fmt.Errorf("pdftk command failed: %w, output: %s", err, string(output))
// 	}

// 	// Read the output file
// 	result, err := os.ReadFile(outputFile)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read output PDF: %w", err)
// 	}

// 	return result, nil
// }
