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
	fmt.Fprintf(w, "こんにちは")
}

// 星座のリスト
var signs = []string{
	"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
	"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
}

// 簡易的な太陽黄経の計算ロジック（近似値）
func calculateSunSign(t time.Time) string {
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
	return signs[index]
}

func handleSunSign(w http.ResponseWriter, r *http.Request) {
    // 1. クエリパラメータ "date" を取得
    dateStr := r.URL.Query().Get("date")
    
    // 2. 緯度・経度もついでに受け取れるようにしておく（将来用）
    latStr := r.URL.Query().Get("lat")
    lonStr := r.URL.Query().Get("lon")

    // 3. 日時のパース。RFC3339（2006-01-02T15:04:05Z）を期待
    t, err := time.Parse(time.RFC3339, dateStr)
    if err != nil {
        // パースに失敗した場合は、今日の日付をデフォルトにする
        t = time.Now()
    }

    sign := calculateSunSign(t)

    // レスポンスを作成
    res := map[string]interface{}{
        "date": t.Format(time.RFC3339),
        "sign": sign,
        "location": map[string]string{
            "lat": latStr,
            "lon": lonStr,
        },
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(res)
}