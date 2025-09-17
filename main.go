package main

import (
	"fmt"
	"log"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// func main() {

// }

// // 1) PDF faylni []byte ko‘rinishida o‘qib olamiz
// inputPDF, err := os.ReadFile("input.pdf")
// if err != nil {
// 	log.Fatalf("PDF o‘qishda xatolik: %v", err)
// }

// // 2) QR code ni []byte ko‘rinishida yaratamiz
// qrBytes, err := qrcode.Encode("https://example.com", qrcode.Medium, 256)
// if err != nil {
// 	log.Fatalf("QR code yaratishda xatolik: %v", err)
// }

// // 3) Biz yozgan test3Bytes funksiyasini chaqiramiz
// outputBytes, err := test3Bytes(inputPDF, qrBytes)
// if err != nil {
// 	log.Fatalf("Xatolik: %v", err)
// }

// // 4) Natijani faylga yozib ko‘ramiz
// err = os.WriteFile("output.pdf", outputBytes, 0644)
// if err != nil {
// 	log.Fatalf("Natijani yozishda xatolik: %v", err)
// }

// 	test2()
// 	fmt.Println("✅ PDF muvaffaqiyatli yaratildi: output.pdf")
// }

// // // PDF oxiriga matn qo'shish funksiyasi
// func addTextToPDF(inputPath, outputPath, text string) error {
// 	// 1-usul: Stamp yordamida matn qo'shish
// 	// return addTextWithStamp(inputPath, outputPath, text)
// 	// 2-usul: Yangi sahifa qo'shib, unga matn yozish
// 	// return addTextWithNewPage(inputPath, outputPath, text)
// 	// 3-usul: Batafsil parametrlar bilan matn qo'shish
// 	// return addTextWithDetailedOptions(inputPath, outputPath, text)
// 	// 4-usul: Ko'p qatorli matn qo'shish
// 	lines := []string{"Birinchi qator", "Ikkinchi qator", "Uchinchi qator"}
// 	return addMultilineText(inputPath, outputPath, lines)
// 	// 5-usul: Maxsus pozitsiya bilan matn qo'shish
// 	// return addTextAtCustomPosition(inputPath, outputPath, text, 100, 200)
// }

// // Stamp yordamida matn qo'shish
// func addTextWithStamp(inputPath, outputPath, text string) error {
// 	// Watermark konfiguratsiyasi
// 	wm, err := pdfcpu.ParseTextWatermarkDetails(text, "font:Helvetica, points:12, pos:bl, off:50 50", false, types.POINTS)
// 	if err != nil {
// 		return fmt.Errorf("watermark yaratishda xatolik: %v", err)
// 	}

// 	// PDF faylni o'qish
// 	conf := model.NewDefaultConfiguration()

// 	// Watermark qo'shish
// 	err = api.AddWatermarksFile(inputPath, outputPath, nil, wm, conf)
// 	if err != nil {
// 		return fmt.Errorf("matn qo'shishda xatolik: %v", err)
// 	}

// 	return nil
// }

// // Yangi sahifa qo'shib, unga matn yozish
// func addTextWithNewPage(inputPath, outputPath, text string) error {
// 	conf := model.NewDefaultConfiguration()

// 	// Avval PDF ni o'qimiz
// 	ctx, err := api.ReadContextFile(inputPath)
// 	if err != nil {
// 		return fmt.Errorf("PDF o'qishda xatolik: %v", err)
// 	}
// 	pageCount := ctx.PageCount // mavjud sahifalar soni

// 	// Oxiridan keyin qo‘shish uchun: n+1
// 	insertAt := []string{fmt.Sprintf("%d", pageCount)} // masalan, "6"

// 	// Yangi sahifa qo'shamiz
// 	err = api.InsertPagesFile(inputPath, "temp.pdf", insertAt, false, pdfcpu.DefaultPageConfiguration(), conf)
// 	if err != nil {
// 		return fmt.Errorf("sahifa qo'shishda xatolik: %v", err)
// 	}

// 	// Yangi sahifaga matn qo'shamiz
// 	wm, err := pdfcpu.ParseTextWatermarkDetails(text, "font:Courier, points:14, pos:c", false, types.POINTS)
// 	if err != nil {
// 		return fmt.Errorf("watermark yaratishda xatolik: %v", err)
// 	}

// 	err = api.AddWatermarksFile("temp.pdf", outputPath, []string{"-1"}, wm, conf)
// 	if err != nil {
// 		return fmt.Errorf("matn qo'shishda xatolik: %v", err)
// 	}

// 	// Vaqtinchalik faylni o'chirish
// 	os.Remove("temp.pdf")

// 	return nil
// }

// // Batafsil parametrlar bilan matn qo'shish
// func addTextWithDetailedOptions(inputPath, outputPath, text string) error {
// 	// Batafsil konfiguratsiya
// 	options := "font:Times-Roman, points:16, pos:bc, rot:0, op:0.8, col:0.2 0.2 0.8, off:0 30"

// 	wm, err := pdfcpu.ParseTextWatermarkDetails(text, options, false)
// 	if err != nil {
// 		return fmt.Errorf("watermark yaratishda xatolik: %v", err)
// 	}

// 	conf := model.NewDefaultConfiguration()

// 	err = api.AddWatermarksFile(inputPath, outputPath, nil, wm, conf)
// 	if err != nil {
// 		return fmt.Errorf("matn qo'shishda xatolik: %v", err)
// 	}

// 	return nil
// }

// // Ko'p qatorli matn qo'shish
// func addMultilineText(inputPath, outputPath string, lines []string) error {
// 	conf := model.NewDefaultConfiguration()

// 	// Har bir qator uchun alohida watermark
// 	for i, line := range lines {
// 		tempOutput := fmt.Sprintf("temp_%d.pdf", i)
// 		currentInput := inputPath

// 		if i > 0 {
// 			currentInput = fmt.Sprintf("temp_%d.pdf", i-1)
// 		}

// 		// Y pozitsiyasini har safar o'zgartirish
// 		yOffset := 50 + (i * 20) // Har qator uchun 20 piksel farq
// 		options := fmt.Sprintf("font:Helvetica, points:10, pos:bl, off:50 %d", yOffset)

// 		wm, err := pdfcpu.ParseTextWatermarkDetails(line, options, false, types.POINTS)
// 		if err != nil {
// 			return fmt.Errorf("watermark yaratishda xatolik: %v", err)
// 		}

// 		finalOutput := outputPath
// 		if i < len(lines)-1 {
// 			finalOutput = tempOutput
// 		}

// 		err = api.AddWatermarksFile(currentInput, finalOutput, nil, wm, conf)
// 		if err != nil {
// 			return fmt.Errorf("matn qo'shishda xatolik: %v", err)
// 		}

// 		// Vaqtinchalik fayllarni tozalash
// 		if i > 0 {
// 			os.Remove(currentInput)
// 		}
// 	}

// 	return nil
// }

// // Maxsus pozitsiya bilan matn qo'shish
// func addTextAtCustomPosition(inputPath, outputPath, text string, x, y float64) error {
// 	options := fmt.Sprintf("font:Helvetica, points:12, pos:tl, off:%.0f %.0f", x, y)

// 	wm, err := pdfcpu.ParseTextWatermarkDetails(text, options, false, types.POINTS)
// 	if err != nil {
// 		return fmt.Errorf("watermark yaratishda xatolik: %v", err)
// 	}

// 	conf := model.NewDefaultConfiguration()

// 	err = api.AddWatermarksFile(inputPath, outputPath, nil, wm, conf)
// 	if err != nil {
// 		return fmt.Errorf("matn qo'shishda xatolik: %v", err)
// 	}

// 	return nil
// }

// func test() {
// Unique abbreviations are accepted for all watermark descriptor parameters.
// eg. sc = scalefactor or rot = rotation

// Add a "Demo" watermark to all pages of in.pdf along the diagonal running from lower left to upper right.
// onTop := false
// update := false
// wm, err := api.TextWatermark("Demo", "", onTop, update, types.POINTS)
// if err != nil {
// 	fmt.Println("Error creating watermark: 167 ", err)
// }
// err = api.AddWatermarksFile("input.pdf", "output.pdf", nil, wm, nil)
// if err != nil {
// 	fmt.Println("Error adding watermark: 171 ", err)
// }

// Stamp all odd pages of input.pdf in red "Confidential" in 48 point Courier
// using a rotation angle of 45 degrees and an absolute scalefactor of 1.0.
// onTop = true
// wm, err = api.TextWatermark("Confidential", "font:Courier, points:48, col: 1 0 0, rot:45, scale:1 abs", onTop, update, types.POINTS)
// if err != nil {
// 	fmt.Println("Error creating watermark: 176 ", err)
// }
// err = api.AddWatermarksFile("input.pdf", "output.pdf", []string{"odd"}, wm, nil)
// if err != nil {
// 	fmt.Println("Error adding watermark: 183 ", err)
// }

// Add image stamps to input.pdf using absolute scaling and a negative rotation of 90 degrees.
// wm, err = api.ImageWatermark("image.png", "scalefactor:.5 a, rot:-90, pos:tl, off:0 0", onTop, update, types.POINTS)
// if err != nil {
// 	fmt.Println("Error creating watermark: ImageWatermark: ", err)
// }
// err = api.AddWatermarksFile("input.pdf", "output.pdf", nil, wm, nil)
// if err != nil {
// 	fmt.Println("Error adding watermark: 193 ", err)
// }

// Add a PDF stamp to all pages of input.pdf using the 2nd page of stamp.pdf, use absolute scaling of 0.5
// and rotate along the 2nd diagonal running from upper left to lower right corner.
// wm, err = api.PDFWatermark("stamp.pdf:2", "scale:.5 abs, diagonal:2", onTop, update, types.POINTS)
// if err != nil {
// 	fmt.Println("Error creating watermark: PDFWatermark: ", err)
// }
// err = api.AddWatermarksFile("input.pdf", "output.pdf", nil, wm, nil)
// if err != nil {
// 	fmt.Println("Error adding watermark: 204 ", err)
// }
// }

func test2() {
	inFile := "input.pdf"
	outFile := "output.pdf"
	imageFile := "qrcode.png"

	// Avval input.pdf ni output.pdf ga ko‘chirib olamiz
	err := api.OptimizeFile(inFile, outFile, nil)
	if err != nil {
		log.Fatalf("Failed to prepare output file: %v", err)
	}

	// PDF ni o‘qib, sahifa sonini aniqlaymiz
	ctx, err := api.ReadContextFile(outFile)
	if err != nil {
		log.Fatalf("PDF o‘qishda xatolik: %v", err)
	}
	pageCount := ctx.PageCount
	targetPage := fmt.Sprintf("%d", pageCount)

	// 1) QR code qo‘shish
	qrWm, err := api.ImageWatermark(imageFile, "pos:tc, scale:0.16, off:-150 -50, rot:0", true, true, types.POINTS)
	if err != nil {
		fmt.Println("Error creating watermark: ImageWatermark: ", err)
	}
	err = api.AddWatermarksFile(outFile, outFile, []string{targetPage}, qrWm, nil)
	if err != nil {
		log.Println("Error adding QR code watermark: ", err)
	}

	// 2) Chapdagi matn
	leftText := model.DefaultWatermarkConfig()
	leftText.Mode = model.WMText
	leftText.TextString = "“Tayyorlandi”\n" +
		"“NAVOIY VILOYATI UCHQUDUQ TUMAN\n" +
		"KAMBAG'ALLIKNI QISQARTIRISH VA\n" +
		"BANDLIKKA KO'MAKLASHISH BO'LIMI”\n" +
		"DAVLAT MUASSASASI\n" +
		"Test Xususiy sherik uchun"
	leftText.Pos = types.TopCenter
	leftText.Dx = -150
	leftText.Dy = -150
	leftText.FontName = "Times-Roman"
	leftText.FontSize = 10
	leftText.ScaledFontSize = 10
	leftText.Scale = 0.4
	leftText.Color = color.Black
	leftText.StrokeColor = color.Black
	leftText.FillColor = color.Black
	leftText.Rotation = 0
	leftText.Diagonal = 0

	err = api.AddWatermarksFile(outFile, outFile, []string{targetPage}, leftText, nil)
	if err != nil {
		log.Println("Error adding left text watermark: ", err)
	}

	rightText := model.DefaultWatermarkConfig()
	rightText.Mode = model.WMText
	rightText.TextString = "“Kelishildi”\n" +
		"“O‘ZBEKISTON RESPUBLIKASI IQTISODIYOT\n" +
		"VA MOLIYA VAZIRLIGI HUZURIDAGI\n" +
		"AXBOROT TEXNOLOGIYALARI MARKAZI”\n" +
		"DAVLAT UNITAR KORXONASI\n" +
		"AXMADOV DILMUROD ELMUROD O‘G‘LI"
	rightText.Pos = types.TopCenter
	rightText.Dx = 150
	rightText.Dy = -50
	rightText.FontName = "Times-Roman"
	rightText.FontSize = 10
	rightText.ScaledFontSize = 10
	rightText.Scale = 0.4
	rightText.Color = color.Black
	rightText.StrokeColor = color.Black
	rightText.FillColor = color.Black
	rightText.Rotation = 0
	rightText.Diagonal = 0

	err = api.AddWatermarksFile(outFile, outFile, []string{targetPage}, rightText, nil)
	if err != nil {
		log.Println("Error adding right text watermark: ", err)
	}

	log.Println("Block added successfully at the end!")
}
