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
			QRContext: "https://example.com/signature2",
			TextList: []string{
				"“Kelishildi”",
				"“O‘ZBEKISTON RESPUBLIKASI IQTISODIYOT VA MOLIYA VAZIRLIGI HUZURIDAGI AXBOROT TEXNOLOGIYALARI MARKAZI” DAVLAT UNITAR KORXONASI",
				"AXMADOV DILMUROD ELMUROD O‘G‘LI",
				"2023 yil 16 fevral",
			},
		},
		// {
		// 	QRContext: "https://example.com/signature2",
		// 	TextList: []string{
		// 		"“Kelishildi”",
		// 		"“O‘ZBEKISTON RESPUBLIKASI IQTISODIYOT VA MOLIYA VAZIRLIGI HUZURIDAGI AXBOROT TEXNOLOGIYALARI MARKAZI” DAVLAT UNITAR KORXONASI",
		// 		"AXMADOV DILMUROD ELMUROD O‘G‘LI",
		// 		"2023 yil 16 fevral",
		// 	},
		// },
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

	pageCount := ctx.PageCount
	wmList := []*model.Watermark{}

	textX := -150.0
	textY := -150.0

	text, _ := textFormatting(signatureList[0].TextList)

	textWm := model.DefaultWatermarkConfig()
	textWm.Mode = model.WMText
	textWm.TextString = text
	textWm.Pos = types.TopCenter
	textWm.Dx = textX
	textWm.Dy = textY
	textWm.FontName = "Times-Roman"
	textWm.FontSize = 10
	textWm.ScaledFontSize = 10
	textWm.Scale = 0.4
	textWm.Color = color.Black
	textWm.StrokeColor = color.Black
	textWm.FillColor = color.Black
	textWm.Rotation = 0
	textWm.Diagonal = 0

	wmList = append(wmList, textWm)

	// add watermarks
	err = api.AddWatermarksSliceMap(in, out, map[int][]*model.Watermark{pageCount: wmList}, nil)
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
