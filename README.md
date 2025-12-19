# Web Scraper 
Bu proje, Siber Tehdit İstihbaratı (CTI) için hazırlanmış bir web scraper uygulamasıdır. Go (Golang) dilinde yazılmış olup, chromedp kütüphanesi kullanarak web sitelerinden HTML içeriği çeker ve ekran görüntüsü alır.

## Hedefleri
✅ Belirtilen URL'den HTML içeriği çekme
✅ Web sayfasının ekran görüntüsünü alma
✅ Hata kontrolü ve kullanıcıya bilgi verme
✅ Çekilen veriyi dosyaya kaydetme
✅ Sayfadaki tüm URL'leri listeleme


### Gereksinimler
- Go 1.21 veya üzeri
- Google Chrome veya Chromium tarayıcısı

### Adım 1: Bağımlılıkları Yükle
```bash
go mod download
```

### Adım 2: Programı Çalıştır
```bash
go run main.go <URL>
```

### Temel Kullanım
```bash
go run main.go https://example.com
```


##  Özellikler
1. **URL'ye Bağlanma:** Chromedp kullanarak gerçek bir tarayıcı ortamında bağlantı kurar
2. **Hata Kontrolü:** Timeout, bağlantı hataları ve diğer sorunları yakalar
3. **Veri Kaydetme:** HTML içeriğini ve screenshot'u yerel dosyalara kaydeder
4. **Link Çıkarma:** Sayfadaki tüm URL'leri toplar ve listeler
