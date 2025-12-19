package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"
	"github.com/chromedp/chromedp"
)
func main() {
	log.SetOutput(io.Discard)

	if len(os.Args) < 2 {
		fmt.Println("Hedef Site URL'si giriniz:")
		os.Exit(1)
	}

	targetURL := os.Args[1]

	fmt.Println("\n--- AŞAMA 1: HTML Bilgisi ---")
	fmt.Println("->Bağlantı kuruluyor:", targetURL)

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var htmlContent string
	var screenshotBuf []byte

	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		fmt.Println("[X] Bağlantı başarısız!")
		os.Exit(1)
	}

	fmt.Println("[+] Bağlantı başarılı (200 OK)")

	htmlFileName := "site.html"
	os.WriteFile(htmlFileName, []byte(htmlContent), 0644)
	fmt.Printf("[+] Sayfanın HTML icerigi '%s' dosyasina kaydedildi.\n", htmlFileName)

	fmt.Println("\n--- AŞAMA 2: Ekran Görüntüsü Alma ---")
	fmt.Println("->Chrome başlatılıyor, lütfen bekleyiniz...")

	chromedp.Run(ctx,
		chromedp.FullScreenshot(&screenshotBuf, 90),
	)

	fmt.Println("-> Siteye gidiliyor ve görüntü işleniyor...")

	screenshotFileName := "screenshot.png"
	os.WriteFile(screenshotFileName, screenshotBuf, 0644)
	fmt.Printf("[+] Ekran görüntüsü başarıyla '%s' dosyasına kaydedildi.\n", screenshotFileName)

	var links []string
	chromedp.Run(ctx,
		chromedp.Evaluate(`Array.from(document.querySelectorAll('a')).map(a => a.href)`, &links),
	)

	if len(links) > 0 {
		linksFileName := "links.txt"
		file, err := os.Create(linksFileName)
		if err == nil {
			defer file.Close()
			file.WriteString(fmt.Sprintf("Toplam %d link bulundu:\n\n", len(links)))
			for i, link := range links {
				file.WriteString(fmt.Sprintf("%d. %s\n", i+1, link))
			}
			fmt.Printf("[+] Sayfadaki linkler '%s' dosyasına kaydedildi. (Toplam: %d)\n", linksFileName, len(links))
		}
	}

	fmt.Println("\n[+] Tüm görevler başarıyla tamamlandı.")
}
