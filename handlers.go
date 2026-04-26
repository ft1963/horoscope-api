package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// --- ハンドラ関数の定義 ---

func handleHello(w http.ResponseWriter, r *http.Request) {
    // 別のパス（例: /anything）が来た時もここが呼ばれるのを防ぐためのチェック
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
	
	html := `<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <title>星座占い・天中殺</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <style>
        body { font-family: sans-serif; max-width: 600px; margin: 0 auto; padding: 2rem; }
        .result-box { margin-top: 1rem; padding: 1rem; border: 1px solid #ccc; border-radius: 4px; background-color: #f9f9f9; }
    </style>
</head>
<body>
    <h1>星座占い・天中殺</h1>
    <h2>星座占い</h2>
    <form hx-post="/htmx-sun-sign" hx-target="#result">
        <label for="dob">生年月日:</label>
        <input type="date" id="dob" name="date" required>
        <button type="submit">鑑定</button>
    </form>
    <div id="result"></div>

    <h2>天中殺</h2>
    <form hx-post="/ten" hx-target="#result-ten">
        <label for="dob-ten">生年月日:</label>
        <input type="date" id="dob-ten" name="date" required>
        <button type="submit">鑑定</button>
    </form>
    <div id="result-ten"></div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func handleHtmxSunSign(w http.ResponseWriter, r *http.Request) {
	dateStr := r.FormValue("date")
	t := parseDate(dateStr)
	signIndex := calculateSunSign(t)
	sign := signs[signIndex]

	html := fmt.Sprintf(`
		<div class="result-box">
			<h2>鑑定結果</h2>
			<p>あなたの星座は <strong>%s</strong> です。</p>
			<p>特徴: %s</p>
		</div>
	`, sign.Name, sign.Trait)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func handleSunSign(w http.ResponseWriter, r *http.Request) {
    // 1. クエリパラメータ "date" を取得
    dateStr := r.URL.Query().Get("date")
    
    // 2. 緯度・経度もついでに受け取れるようにしておく（将来用）
    latStr := r.URL.Query().Get("lat")
    lonStr := r.URL.Query().Get("lon")

    // 3. 日時のパース（共通処理化）
    t := parseDate(dateStr)

    signIndex := calculateSunSign(t)

    // レスポンスを作成
    res := map[string]interface{}{
        "date":  t.Format("2006-01-02T15:04:05Z07:00"), // time.RFC3339
        "sign":  signs[signIndex].Name,
        "trait": signs[signIndex].Trait,
        "location": map[string]string{
            "lat": latStr,
            "lon": lonStr,
        },
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(res)
}

// HTMX用の天中殺ハンドラ
func handleTen(w http.ResponseWriter, r *http.Request) {
	dateStr := r.FormValue("date")
	t := parseDate(dateStr)
	
	y, m, d := t.Year(), int(t.Month()), t.Day()
	
	dayK := getDayEto(y, m, d)
	tenchu := getTenchusatsu(dayK)

	html := fmt.Sprintf(`
		<div class="result-box">
			<h2>天中殺 鑑定結果</h2>
			<p>あなたの日干支は <strong>%s</strong> です。</p>
			<p>天中殺: <strong>%s</strong></p>
		</div>
	`, dayK, tenchu)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}
