package main

import "math"

// --- 天中殺の計算 ---

// 干支の配列を定義
func getEtoArrays() ([]string, []string) {
	kan := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	shi := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	return kan, shi
}

// ユリウス日を計算
func julianDayFromDate(y, m, d int) int {
	if m <= 2 {
		y--
		m += 12
	}
	a := y / 100
	b := 2 - a + (a / 4)
	
	// PHPのfloor(365.25 * ...)等の浮動小数点演算を整数演算で再現
	term1 := int(math.Floor(365.25 * float64(y+4716)))
	term2 := int(math.Floor(30.6001 * float64(m+1)))
	
	return term1 + term2 + d + b - 1524
}

// 日干支を取得
func getDayEto(y, m, d int) string {
	if y == 0 || m == 0 || d == 0 {
		return ""
	}

	kan, shi := getEtoArrays()

	jd := julianDayFromDate(y, m, d)
	base := julianDayFromDate(1920, 1, 7) // 1920年1月7日は「甲子」

	diff := jd - base

	// 負の数に対応した剰余計算
	kIdx := (diff%10 + 10) % 10
	sIdx := (diff%12 + 12) % 12

	return kan[kIdx] + shi[sIdx]
}

// 干支の番号（1-60）を取得
func getEtoNumber(eto string) int {
	kan, shi := getEtoArrays()

	for i := 0; i < 60; i++ {
		if kan[i%10]+shi[i%12] == eto {
			return i + 1
		}
	}
	return 0
}

// 天中殺を取得
func getTenchusatsu(dayEto string) string {
	num := getEtoNumber(dayEto)
	if num == 0 {
		return ""
	}

	groups := []string{
		"戌亥天中殺",
		"申酉天中殺",
		"午未天中殺",
		"辰巳天中殺",
		"寅卯天中殺",
		"子丑天中殺",
	}

	// (num-1)/10 は整数除算なので PHPの intval と同等
	return groups[(num-1)/10]
}
