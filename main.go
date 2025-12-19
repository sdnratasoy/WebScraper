package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
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

	// HTTP durum kodu kontrolü
	resp, err := http.Get(targetURL)
	if err != nil {
		fmt.Println("[X] Bağlantı hatası:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Println("[X] Hata: 404 Sayfa Bulunamadı")
		fmt.Println("Belirtilen sayfa mevcut değil.")
		os.Exit(1)
	} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		fmt.Printf("[X] Hata: %d İstemci Hatası\n", resp.StatusCode)
		fmt.Println("İsteğinizde bir sorun var.")
		os.Exit(1)
	} else if resp.StatusCode >= 500 {
		fmt.Printf("[X] Hata: %d Sunucu Hatası\n", resp.StatusCode)
		fmt.Println("Sunucu yanıt vermiyor veya bir hata oluştu.")
		os.Exit(1)
	}

	fmt.Printf("[+] Bağlantı başarılı (%d %s)\n", resp.StatusCode, http.StatusText(resp.StatusCode))

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var htmlContent string
	var screenshotBuf []byte
	
	err = chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		fmt.Println("[X] Sayfa yüklenirken hata oluştu:", err)
		os.Exit(1)
	}

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
