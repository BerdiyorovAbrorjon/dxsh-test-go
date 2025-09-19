package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/skip2/go-qrcode"
)

const (
	maxLineLength = 32
)

// Rectangle represents a rectangular area in PDF coordinates
type Rectangle struct {
	X      float64 // Left edge
	Y      float64 // Bottom edge (PDF coordinate system)
	Width  float64
	Height float64
}

// SignaturePlacement represents where a signature should be placed
type SignaturePlacement struct {
	Signature SignatureInfo
	QRRect    Rectangle
	TextRect  Rectangle
	Page      int
}

type SignatureInfo struct {
	QRContext string
	TextList  []string
}

func main() {
	inFile := "content.pdf" // Use PDF with actual content for testing
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

	err = os.WriteFile("sample.pdf", out, 0644)
	if err != nil {
		log.Fatalf("Natija PDF faylini yozishda xatolik: %v", err)
	}

	log.Println("PDF ga QR code muvaffaqiyatli qoshildi: sample.pdf")
}

func AddSignatureListToPDF(pdfData []byte, signatureList []SignatureInfo) ([]byte, error) {
	in := bytes.NewReader(pdfData)
	out := new(bytes.Buffer)

	ctx, err := api.ReadAndValidate(in, model.NewDefaultConfiguration())
	if err != nil {
		return nil, fmt.Errorf("PDF oqishda xatolik: %v", err)
	}
	originalPageCount := ctx.PageCount

	// Step 1: Always create new pages for ALL signatures (no placement on existing pages)
	var signaturePlacements []SignaturePlacement

	log.Printf("Creating new pages for all %d signatures", len(signatureList))

	// Step 2: Calculate how many new pages we need
	signaturesPerNewPage := 6 // 3x2 grid on new pages
	newPagesNeeded := (len(signatureList) + signaturesPerNewPage - 1) / signaturesPerNewPage

	// Step 3: Add new pages
	currentData := pdfData
	newPagesAdded := 0

	for i := 0; i < newPagesNeeded; i++ {
		in = bytes.NewReader(currentData)
		out.Reset()

		err = api.InsertPages(in, out, []string{fmt.Sprintf("%d", originalPageCount+i)}, false, pdfcpu.DefaultPageConfiguration(), nil)
		if err != nil {
			return nil, fmt.Errorf("sahifa qo'shishda xatolik: %v", err)
		}
		currentData = out.Bytes()
		newPagesAdded++
	}

	log.Printf("Added %d new pages for all signatures", newPagesAdded)

	// Step 4: Place ALL signatures on new pages with grid layout
	for i, signature := range signatureList {
		pageIndex := i / signaturesPerNewPage
		positionOnPage := i % signaturesPerNewPage
		actualPage := originalPageCount + pageIndex + 1

		// Grid layout for new pages
		row := positionOnPage / 2
		col := positionOnPage % 2

		// Calculate positions using ORIGINAL PERFECT system for new pages
		startX := -150.0 + float64(col)*300.0
		startY := -120.0 - float64(row)*280.0

		qrRect := Rectangle{
			X: startX, Y: startY + 90, // QR above text - ORIGINAL PERFECT SPACING
			Width: 60, Height: 60, // ORIGINAL PERFECT SIZE
		}

		textRect := Rectangle{
			X: startX, Y: startY,
			Width: 150, Height: 80, // ORIGINAL PERFECT SIZE
		}

		placement := SignaturePlacement{
			Signature: signature,
			QRRect:    qrRect,
			TextRect:  textRect,
			Page:      actualPage,
		}

		signaturePlacements = append(signaturePlacements, placement)

		log.Printf("Placed signature %d on new page %d (row %d, col %d)",
			i+1, actualPage, row, col)
	}

	// Step 8: Convert placements to watermarks and apply them
	pageWatermarks := make(map[int][]*model.Watermark)

	for _, placement := range signaturePlacements {
		// Generate QR code
		qrBytes, err := qrCodeGenerate(placement.Signature.QRContext)
		if err != nil {
			return nil, fmt.Errorf("QR code yaratishda xatolik: %v", err)
		}

		// Create QR watermark
		qrWm := model.DefaultWatermarkConfig()
		qrWm.Mode = model.WMImage
		qrWm.Image = bytes.NewReader(qrBytes)

		// Use TopCenter positioning for ALL pages - SAME SYSTEM
		qrWm.Pos = types.TopCenter
		qrWm.Dx = placement.QRRect.X
		qrWm.Dy = placement.QRRect.Y

		qrWm.Scale = 0.15 // ORIGINAL PERFECT SCALE
		qrWm.Rotation = 0
		qrWm.Diagonal = 0
		qrWm.Opacity = 0.95 // Keep only transparency to reduce white background

		// Create text watermark
		text, _ := textFormatting(placement.Signature.TextList)

		textWm := model.DefaultWatermarkConfig()
		textWm.Mode = model.WMText
		textWm.TextString = text

		// Use TopCenter positioning for ALL pages - SAME SYSTEM
		textWm.Pos = types.TopCenter
		textWm.Dx = placement.TextRect.X
		textWm.Dy = placement.TextRect.Y

		textWm.FontName = "Times-Roman"
		textWm.FontSize = 9       // ORIGINAL PERFECT SIZE
		textWm.ScaledFontSize = 9 // ORIGINAL PERFECT SIZE
		textWm.Scale = 0.35       // ORIGINAL PERFECT SCALE
		textWm.Color = color.Black
		// Try removing stroke and fill colors that might create white backgrounds
		textWm.Rotation = 0
		textWm.Diagonal = 0
		textWm.Opacity = 0.95 // Slight transparency to reduce white background

		// Add to page watermarks
		if pageWatermarks[placement.Page] == nil {
			pageWatermarks[placement.Page] = []*model.Watermark{}
		}
		pageWatermarks[placement.Page] = append(pageWatermarks[placement.Page], qrWm, textWm)
	}

	// Step 9: Apply all watermarks
	in = bytes.NewReader(currentData)
	out.Reset()
	err = api.AddWatermarksSliceMap(in, out, pageWatermarks, nil)
	if err != nil {
		return nil, fmt.Errorf("error on add watermarks to PDF: %v", err)
	}

	log.Printf("Successfully placed %d signatures on %d new pages",
		len(signaturePlacements),
		newPagesAdded)

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

// min returns the minimum of two int values
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
