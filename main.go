package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/skip2/go-qrcode"
)

const (
	maxLineLength = 32
)

type SignatureInfo struct {
	QRContext string
	TextList  []string
}

func main() {

	inFile := "blank.pdf"
	signatureList := []SignatureInfo{
		{
			QRContext: "https://example.com/signature1",
			TextList: []string{
				"“Tayyorlandi”",
				"“NAVOIY VILOYATI UCHQUDUQ TUMAN KAMBAG'ALLIKNI QISQARTIRISH VA BANDLIKKA KO'MAKLASHISH BO'LIMI” DAVLAT MUASSASASI",
				"KARIMOV JAHONGIR RAXMONBERDI O‘G‘LI",
				"2023 yil 15 fevral",
			},
		},
		{
			QRContext: "https://example.com/signature1",
			TextList: []string{
				"“Tayyorlandi”",
				"“NAVOIY VILOYATI UCHQUDUQ TUMAN KAMBAG'ALLIKNI QISQARTIRISH VA BANDLIKKA KO'MAKLASHISH BO'LIMI” DAVLAT MUASSASASI",
				"KARIMOV JAHONGIR RAXMONBERDI O‘G‘LI",
				"2023 yil 15 fevral",
			},
		},
		{
			QRContext: "https://example.com/signature2",
			TextList: []string{
				"“Kelishildi”",
				"“O‘ZBEKISTON RESPUBLIKASI IQTISODIYOT VA MOLIYA VAZIRLIGI HUZURIDAGI AXBOROT TEXNOLOGIYALARI MARKAZI” DAVLAT UNITAR KORXONASI",
				"AXMADOV DILMUROD ELMUROD O‘G‘LI",
				"2023 yil 16 fevral",
			},
		},
		{
			QRContext: "https://example.com/signature2",
			TextList: []string{
				"“Kelishildi”",
				"“O‘ZBEKISTON RESPUBLIKASI IQTISODIYOT VA MOLIYA VAZIRLIGI HUZURIDAGI AXBOROT TEXNOLOGIYALARI MARKAZI” DAVLAT UNITAR KORXONASI",
				"AXMADOV DILMUROD ELMUROD O‘G‘LI",
				"2023 yil 16 fevral adsfasd afdasdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf ",
			},
		},
		{
			QRContext: "https://example.com/signature2",
			TextList: []string{
				"“Kelishildi”",
				"“O‘ZBEKISTON RESPUBLIKASI IQTISODIYOT VA MOLIYA VAZIRLIGI HUZURIDAGI AXBOROT TEXNOLOGIYALARI MARKAZI” DAVLAT UNITAR KORXONASI",
				"AXMADOV DILMUROD ELMUROD O‘G‘LI",
				"2023 yil 16 fevral",
			},
		},
	}

	pdfFile, err := os.Open(inFile)
	if err != nil {
		panic(err)
	}
	defer pdfFile.Close()

	pdfData, err := io.ReadAll(pdfFile)
	if err != nil {
		panic(err)
	}

	out, err := AddSignatureListToPDF(pdfData, signatureList)
	if err != nil {
		log.Fatalf("PDF ga QR code qoshishda xatolik: %v", err)
	}

	err = os.WriteFile("output.pdf", out, 0644)
	if err != nil {
		log.Fatalf("Natija PDF faylini yozishda xatolik: %v", err)
	}

	log.Println("PDF ga QR code muvaffaqiyatli qoshildi:output.pdf")
}

func AddSignatureListToPDF(pdfData []byte, signatureList []SignatureInfo) ([]byte, error) {
	in := bytes.NewReader(pdfData)
	out := new(bytes.Buffer)

	ctx, err := api.ReadAndValidate(in, model.NewDefaultConfiguration())
	if err != nil {
		return nil, fmt.Errorf("PDF oqishda xatolik: %v", err)
	}

	var pageWidth = 595.0  // A4 width (210 mm)
	var pageHeight = 842.0 // A4 height (297 mm)

	// Simple and working layout configuration
	itemsPerRow := 2           // 2 signatures per row
	
	// Simple positioning - use negative values for TopCenter positioning
	// TopCenter means offset from center of page
	startX := -150.0           // Start position for first signature
	startY := -120.0           // Start position with proper top margin
	
	horizontalSpacing := 300.0 // Horizontal spacing between signatures
	verticalSpacing := 280.0   // Vertical spacing between rows
	qrTextSpacing := 80.0      // Space between QR and text (QR above text)
	
	// Signature component sizing
	qrSize := 0.15             // QR code scale 
	textScale := 0.35          // Text scale for readability
	fontSize := 9              // Font size for better fit (int type)

	// Simple calculation for max rows
	maxRowsPerPage := 2  // Keep it simple - 2 rows per page
	
	log.Printf("A4 Layout: %.0fx%.0f points, Simple positioning with %d max rows", 
		pageWidth, pageHeight, maxRowsPerPage)

	// Group signatures by pages
	pageWatermarks := make(map[int][]*model.Watermark)

	for i, signature := range signatureList {
		// Simple approach: put all signatures on the last page but with proper positioning
		actualPage := ctx.PageCount

		// Calculate position for proper grid layout
		row := i / itemsPerRow
		col := i % itemsPerRow

		// Calculate position for this signature using optimized spacing
		textX := startX + float64(col)*horizontalSpacing
		textY := startY - float64(row)*verticalSpacing
		qrX := textX
		qrY := textY + qrTextSpacing

		// Generate and add QR code watermark
		qrBytes, err := qrCodeGenerate(signature.QRContext)
		if err != nil {
			return nil, fmt.Errorf("QR code yaratishda xatolik: %v", err)
		}

		qrWm := model.DefaultWatermarkConfig()
		qrWm.Mode = model.WMImage
		qrWm.Image = bytes.NewReader(qrBytes)
		qrWm.Pos = types.TopCenter
		qrWm.Dx = qrX
		qrWm.Dy = qrY
		qrWm.Scale = qrSize  // Use optimized QR size
		qrWm.Rotation = 0
		qrWm.Diagonal = 0

		// Generate and add text watermark
		text, _ := textFormatting(signature.TextList)

		textWm := model.DefaultWatermarkConfig()
		textWm.Mode = model.WMText
		textWm.TextString = text
		textWm.Pos = types.TopCenter
		textWm.Dx = textX
		textWm.Dy = textY
		textWm.FontName = "Times-Roman"
		textWm.FontSize = fontSize         // Use optimized font size
		textWm.ScaledFontSize = fontSize   // Use optimized font size
		textWm.Scale = textScale           // Use optimized text scale
		textWm.Color = color.Black
		textWm.StrokeColor = color.Black
		textWm.FillColor = color.Black
		textWm.Rotation = 0
		textWm.Diagonal = 0

		// Add watermarks to the appropriate page
		if pageWatermarks[actualPage] == nil {
			pageWatermarks[actualPage] = []*model.Watermark{}
		}
		pageWatermarks[actualPage] = append(pageWatermarks[actualPage], qrWm, textWm)
	}

	// For now, we'll put overflow signatures on the last available page
	// This is a simplified approach - in production you'd want to add actual new pages

	// Add watermarks to all pages
	err = api.AddWatermarksSliceMap(in, out, pageWatermarks, nil)
	if err != nil {
		return nil, fmt.Errorf("error on add watermarks to PDF: %v", err)
	}

	return out.Bytes(), nil
}

func qrCodeGenerate(qrContent string) ([]byte, error) {
	qrBytes, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("error on generate QR code: %v", err)
	}
	return qrBytes, nil
}

func textFormatting(textList []string) (string, int) {
	var formattedLines []string

	for _, text := range textList {
		words := strings.Fields(text)
		var currentLine string

		for _, word := range words {
			if len(currentLine)+len(word)+1 > maxLineLength {
				formattedLines = append(formattedLines, currentLine)
				currentLine = word
			} else {
				if currentLine != "" {
					currentLine += " "
				}
				currentLine += word
			}
		}
		if currentLine != "" {
			formattedLines = append(formattedLines, currentLine)
		}
	}

	return strings.Join(formattedLines, "\n"), len(formattedLines)
}
