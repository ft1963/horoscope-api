package main

import (
	"fmt"
	"net/http"
	"os"
	"encoding/json"
	"math"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 各エンドポイントにハンドラ関数を登録
	http.HandleFunc("/", handleHello)
	http.HandleFunc("/sun-sign", handleSunSign)
	http.HandleFunc("/htmx-sun-sign", handleHtmxSunSign)

	fmt.Printf("Server starting on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}

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
    <title>星座占い</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <style>
        body { font-family: sans-serif; max-width: 600px; margin: 0 auto; padding: 2rem; }
        .result-box { margin-top: 1rem; padding: 1rem; border: 1px solid #ccc; border-radius: 4px; background-color: #f9f9f9; }
    </style>
</head>
<body>
    <h1>星座占い</h1>
    <form hx-post="/htmx-sun-sign" hx-target="#result">
        <label for="dob">生年月日:</label>
        <input type="date" id="dob" name="date" required>
        <button type="submit">鑑定</button>
    </form>
    <div id="result"></div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// 星座の情報を持つ構造体
type Sign struct {
	Name  string
	Trait string
}

// 星座のリスト
var signs = []Sign{
	{"牡羊座", "情熱的で行動力がある"},
	{"牡牛座", "マイペースで忍耐強い"},
	{"双子座", "好奇心旺盛でコミュニケーション能力が高い"},
	{"蟹座", "感受性が豊かで家族思い"},
	{"獅子座", "自信に満ちたリーダー気質"},
	{"乙女座", "几帳面で分析力がある"},
	{"天秤座", "社交的でバランス感覚に優れる"},
	{"蠍座", "探究心が強く情熱を秘める"},
	{"射手座", "自由を愛する楽天家"},
	{"山羊座", "真面目で責任感が強い"},
	{"水瓶座", "独創的でマイペース"},
	{"魚座", "共感力が高くロマンチスト"},
}

// 簡易的な太陽黄経の計算ロジック（近似値）
func calculateSunSign(t time.Time) int {
	// 1月1日からの経過日数を計算
	dayOfYear := float64(t.YearDay())
	
	// 春分の日（約3月21日）を0度とした近似計算
	// 実際にはもっと複雑な数式（ユリウス日を使用）が必要ですが、まずは近似値で。
	// 365日で360度回転するので、1日あたり約0.986度
	offset := dayOfYear - 80 // 3月20日頃を基準にする
	if offset < 0 {
		offset += 365
	}
	
	longitude := offset * (360.0 / 365.25)
	index := int(math.Floor(longitude / 30.0))
	
	if index >= 12 {
		index = 0
	}
	return index
}

// 日付文字列をパースする共通関数
func parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Now()
	}
	// HTML5 date input format (YYYY-MM-DD)
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t
	}
	// RFC3339 format
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t
	}
	return time.Now()
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
        "date":  t.Format(time.RFC3339),
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