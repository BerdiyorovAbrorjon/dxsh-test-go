package main

import (
	"fmt"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/skip2/go-qrcode"
)

func main() {
	inputPath := "input.pdf"
	outputPath := "output.pdf"

	// Imzolar ro'yxati (masalan, sign info)
	signatures := []string{"Sign1 info", "Sign2 info", "Sign3 info", "Sign4 info"} // Ko'p bo'lsa

	// PDF faylni yuklash
	file, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// PDF kontekstini yaratish
	ctx, err := api.ReadContext(file, nil)
	if err != nil {
		panic(err)
	}

	// Oxirgi sahifaga o'tish
	numPages := ctx.PageCount
	currentPage := numPages

	// Layout parametrlari
	rowY := 700.0               // A4 sahifa balandligi ~842
	colX1, colX2 := 50.0, 300.0 // Chap va o'ng pozitsiyalar
	itemsPerRow := 2
	rowHeight := 100.0

	for i := 0; i < len(signatures); i += itemsPerRow {
		rowItems := signatures[i:min(i+itemsPerRow, len(signatures))]

		// Agar sahifaga sig'masa, yangi sahifa qo'shish
		if rowY < 50 {
			currentPage++
			rowY = 700.0
			// Yangi kontekstni qayta yuklash
			file, _ := os.Open(outputPath)
			ctx, _ = api.ReadContext(file, nil)
			file.Close()
		}

		// Har qator uchun 2 ta element
		for j, info := range rowItems {
			x := colX1
			if j == 1 {
				x = colX2
			}

			// QR kod generatsiya qilish
			qr, _ := qrcode.New(info, qrcode.Medium)
			qrPath := fmt.Sprintf("qr_%d.png", i+j)
			qr.WriteFile(256, qrPath)

			// QR kodni PDF ga qo'shish
			api.AddImageToPage(inputPath, outputPath, currentPage, qrPath, &pdfcpu.ImageOptions{
				Position: [2]float64{x, rowY - 50},
				Scale:    0.5, // QR o'lchami
			})

			// Matn (sign info) qo'shish
			api.AddText(inputPath, outputPath, currentPage, info, &pdfcpu.TextOptions{
				Position: [2]float64{x, rowY - 80},
				Font:     "Helvetica",
				FontSize: 12,
			})
		}

		rowY -= rowHeight
	}

	// Faylni saqlash
	api.WriteContextFile(ctx, outputPath)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
