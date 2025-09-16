package main

// import (
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"log"
// 	"os"
// 	"strings"

// 	"github.com/pdfcpu/pdfcpu/pkg/api"
// 	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
// 	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
// 	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
// 	"github.com/skip2/go-qrcode"
// )

// const (
// 	maxLineLength = 32
// )

// type SignatureInfo struct {
// 	Order     int
// 	QRContext string
// 	TextList  []string
// }

// func main() {

// 	inFile := "blank.pdf"
// 	signatureList := []SignatureInfo{
// 		{
// 			Order:     1,
// 			QRContext: "https://example.com/signature1",
// 			TextList: []string{
// 				"“Tayyorlandi”",
// 				"“NAVOIY VILOYATI UCHQUDUQ TUMAN KAMBAG'ALLIKNI QISQARTIRISH VA BANDLIKKA KO'MAKLASHISH BO'LIMI” DAVLAT MUASSASASI",
// 				"KARIMOV JAHONGIR RAXMONBERDI O‘G‘LI",
// 				"2023 yil 15 fevral",
// 			},
// 		},
// 		{
// 			Order:     2,
// 			QRContext: "https://example.com/signature2",
// 			TextList: []string{
// 				"“Kelishildi”",
// 				"“O‘ZBEKISTON RESPUBLIKASI IQTISODIYOT VA MOLIYA VAZIRLIGI HUZURIDAGI AXBOROT TEXNOLOGIYALARI MARKAZI” DAVLAT UNITAR KORXONASI",
// 				"AXMADOV DILMUROD ELMUROD O‘G‘LI",
// 				"2023 yil 16 fevral",
// 			},
// 		},
// 		// {
// 		// 	Order:     3,
// 		// 	QRContext: "https://example.com/signature2",
// 		// 	TextList: []string{
// 		// 		"“Kelishildi”",
// 		// 		"“O‘ZBEKISTON RESPUBLIKASI IQTISODIYOT VA MOLIYA VAZIRLIGI HUZURIDAGI AXBOROT TEXNOLOGIYALARI MARKAZI” DAVLAT UNITAR KORXONASI",
// 		// 		"AXMADOV DILMUROD ELMUROD O‘G‘LI",
// 		// 		"2023 yil 16 fevral",
// 		// 	},
// 		// },
// 	}

// 	pdfFile, err := os.Open(inFile)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer pdfFile.Close()

// 	pdfData, err := io.ReadAll(pdfFile)
// 	if err != nil {
// 		panic(err)
// 	}

// 	out, err := AddSignatureListToPDF2(pdfData, signatureList)
// 	if err != nil {
// 		log.Fatalf("PDF ga QR code qoshishda xatolik: %v", err)
// 	}

// 	err = os.WriteFile("output.pdf", out, 0644)
// 	if err != nil {
// 		log.Fatalf("Natija PDF faylini yozishda xatolik: %v", err)
// 	}

// 	log.Println("PDF ga QR code muvaffaqiyatli qoshildi:output.pdf")
// }

// func AddSignatureListToPDF2(pdfData []byte, signatureList []SignatureInfo) ([]byte, error) {
// 	in := bytes.NewReader(pdfData)
// 	out := new(bytes.Buffer)
// 	ctx, err := api.ReadValidateAndOptimize(in, model.NewDefaultConfiguration())
// 	if err != nil {
// 		return nil, fmt.Errorf("PDF oqishda xatolik: %v", err)
// 	}

// 	pageCount := ctx.PageCount
// 	row := 0
// 	col := 0

// 	wmList := []*model.Watermark{}

// 	for _, sig := range signatureList {
// 		qrData, err := qrCodeGenerate(sig.QRContext)
// 		if err != nil {
// 			return nil, fmt.Errorf("error on generate QR code: %v", err)
// 		}
// 		qrReader := bytes.NewBuffer(qrData)

// 		// QR code
// 		qrX := -150 + col*300
// 		qrY := -50 - row*200
// 		qrWm, err := api.ImageWatermarkForReader(
// 			qrReader,
// 			fmt.Sprintf("pos:tc, scale:0.16, off:%d %d, rot:0", qrX, qrY),
// 			true, true, types.POINTS,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("error creating qr watermark: %v", err)
// 		}
// 		wmList = append(wmList, qrWm)

// 		col++
// 		if col > 1 {
// 			col = 0
// 			row++
// 		}
// 	}

// 	err = api.AddWatermarksSliceMap(in, out, map[int][]*model.Watermark{pageCount: wmList}, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("error on add watermarks to PDF: %v", err)
// 	}

// 	return out.Bytes(), nil
// }

// func addSignatureWatermark(pdfBytes, qrBytes []byte, text string, pageNumber, row, col, prevRowLineCount int) ([]byte, error) {
// 	in := bytes.NewReader(pdfBytes)
// 	out := new(bytes.Buffer)
// 	qrReader := bytes.NewBuffer(qrBytes)
// 	targetPage := fmt.Sprintf("%d", pageNumber)

// 	// QR code
// 	qrX := -150 + col*300
// 	qrY := -50 - row*200
// 	qrWm, err := api.ImageWatermarkForReader(qrReader, fmt.Sprintf("pos:tc, scale:0.16, off:%d %d, rot:0", qrX, qrY), true, true, types.POINTS)
// 	if err != nil {
// 		return nil, fmt.Errorf("Error creating qr watermark: %v", err)
// 	}
// 	err = api.AddWatermarks(in, out, []string{targetPage}, qrWm, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("Error on add QR code to PDF: %v", err)
// 	}

// 	// oraliq natija olish
// 	middleBytes := out.Bytes()
// 	in = bytes.NewReader(middleBytes)
// 	out.Reset()

// 	// Text block
// 	textX := -150 + float64(col)*300
// 	textY := -150 - float64(row)*200

// 	textWm := model.DefaultWatermarkConfig()
// 	textWm.Mode = model.WMText
// 	textWm.TextString = text
// 	textWm.Pos = types.TopCenter
// 	textWm.Dx = textX
// 	textWm.Dy = textY
// 	textWm.FontName = "Times-Roman"
// 	textWm.FontSize = 10
// 	textWm.ScaledFontSize = 10
// 	textWm.Scale = 0.4
// 	textWm.Color = color.Black
// 	textWm.StrokeColor = color.Black
// 	textWm.FillColor = color.Black
// 	textWm.Rotation = 0
// 	textWm.Diagonal = 0

// 	err = api.AddWatermarks(in, out, []string{targetPage}, textWm, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("Error on add text to PDF: %v", err)
// 	}

// 	return out.Bytes(), nil
// }

// func qrCodeGenerate(qrContent string) ([]byte, error) {
// 	qrBytes, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
// 	if err != nil {
// 		return nil, fmt.Errorf("Error on generate QR code: %v", err)
// 	}
// 	return qrBytes, nil
// }

// func textFormatting(textList []string) (string, int) {
// 	var formattedLines []string

// 	for _, text := range textList {
// 		words := strings.Fields(text)
// 		var currentLine string

// 		for _, word := range words {
// 			if len(currentLine)+len(word)+1 > maxLineLength {
// 				formattedLines = append(formattedLines, currentLine)
// 				currentLine = word
// 			} else {
// 				if currentLine != "" {
// 					currentLine += " "
// 				}
// 				currentLine += word
// 			}
// 		}
// 		if currentLine != "" {
// 			formattedLines = append(formattedLines, currentLine)
// 		}
// 	}

// 	return strings.Join(formattedLines, "\n"), len(formattedLines)
// }

// func AddSignatureListToPDF(pdfData []byte, signatureList []SignatureInfo) ([]byte, error) {
// 	in := bytes.NewReader(pdfData)
// 	out := new(bytes.Buffer)
// 	middleBytes := out.Bytes()
// 	ctx, err := api.ReadAndValidate(in, nil)
// 	if err != nil {
// 		log.Fatalf("PDF o‘qishda xatolik: %v", err)
// 	}

// 	pageCount := ctx.PageCount
// 	row := 0
// 	col := 0
// 	// prevRowLineCount := 0
// 	currentRowMaxLineCount := 0
// 	targetPage := fmt.Sprintf("%d", pageCount)

// 	// Har bir signature uchun
// 	for _, sig := range signatureList {
// 		qrData, err := qrCodeGenerate(sig.QRContext)
// 		if err != nil {
// 			return nil, fmt.Errorf("error on generate QR code: %v", err)
// 		}
// 		qrReader := bytes.NewBuffer(qrData)

// 		text, lineCount := textFormatting(sig.TextList)
// 		if lineCount > currentRowMaxLineCount {
// 			currentRowMaxLineCount = lineCount
// 		}

// 		// QR code
// 		qrX := -150 + col*300
// 		qrY := -50 - row*200
// 		qrWm, err := api.ImageWatermarkForReader(qrReader, fmt.Sprintf("pos:tc, scale:0.16, off:%d %d, rot:0", qrX, qrY), true, true, types.POINTS)
// 		if err != nil {
// 			return nil, fmt.Errorf("Error creating qr watermark: %v", err)
// 		}
// 		err = api.AddWatermarks(in, out, []string{targetPage}, qrWm, nil)
// 		if err != nil {
// 			return nil, fmt.Errorf("Error on add QR code to PDF: %v", err)
// 		}

// 		// oraliq natija olish
// 		middleBytes = out.Bytes()
// 		in = bytes.NewReader(middleBytes)
// 		out.Reset()

// 		// Text block
// 		textX := -150 + float64(col)*300
// 		textY := -150 - float64(row)*200

// 		textWm := model.DefaultWatermarkConfig()
// 		textWm.Mode = model.WMText
// 		textWm.TextString = text
// 		textWm.Pos = types.TopCenter
// 		textWm.Dx = textX
// 		textWm.Dy = textY
// 		textWm.FontName = "Times-Roman"
// 		textWm.FontSize = 10
// 		textWm.ScaledFontSize = 10
// 		textWm.Scale = 0.4
// 		textWm.Color = color.Black
// 		textWm.StrokeColor = color.Black
// 		textWm.FillColor = color.Black
// 		textWm.Rotation = 0
// 		textWm.Diagonal = 0

// 		err = api.AddWatermarks(in, out, []string{targetPage}, textWm, nil)
// 		if err != nil {
// 			return nil, fmt.Errorf("Error on add text to PDF: %v", err)
// 		}

// 		middleBytes = out.Bytes()
// 		in = bytes.NewReader(middleBytes)
// 		out.Reset()

// 		col++
// 		if col > 1 {
// 			col = 0
// 			row++
// 			// prevRowLineCount = currentRowMaxLineCount
// 			currentRowMaxLineCount = 0
// 		}
// 	}

// 	return middleBytes, nil
// }
