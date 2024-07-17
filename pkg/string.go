package pkg

import (
	"fmt"
	"strconv"
	"strings"
)

// SqlQueryEmail 產生SQL查詢用的字串
//
//	主要目的為避免用戶使用 "+" 重複註冊
//	ex: example@example.com -> example%@example.com
//	ex: example+test@example.com -> example%@example.com
func SqlQueryEmail(email string) string {
	username, domain := EmailSplit(email)
	return fmt.Sprintf("%s%%@%s", username, domain)
}

// EmailSplit 將 email 字串切割成 username,domain 並返回
func EmailSplit(email string) (string, string) {
	emailSplit := strings.Split(email, "@")

	if len(emailSplit) < 2 {
		return "", ""
	}

	var domain string

	domain = emailSplit[1]

	username := strings.Split(emailSplit[0], "+")[0]

	return username, domain
}

// GetRealEmail 取得解析後的真實信箱
//
//	ex: example@example.com -> example@example.com
//	ex: example+test@example.com -> example@example.com
func GetRealEmail(email string) string {
	username, domain := EmailSplit(email)
	return fmt.Sprintf("%s@%s", username, domain)
}

// IsGUINumber 檢查統一編號格式
func IsGUINumber(GUINumber string) bool {
	// 先檢查字串長度
	if len(GUINumber) != 8 {
		return false
	}

	var GUINumberInt []int

	for i := 0; i < 8; i++ {
		n := GUINumber[i]
		nInt, err := strconv.Atoi(string(n))
		if err != nil {
			return false
		}
		GUINumberInt = append(GUINumberInt, nInt)
	}

	var sum int

	vNumber := []int{1, 2, 1, 2, 1, 2, 4, 1}
	for i := 0; i < 8; i++ {
		sum += func() int {
			num := GUINumberInt[i] * vNumber[i]
			// 若 num 大於九 則將 十位數 與 個位數 相加後 回傳
			if num > 9 {
				return num%10 + num/10
			}
			return num
		}()
	}

	switch {
	// 若總和整除 10 返回 true
	case sum%10 == 0:
		return true
	// 若第 7 位 為 "7" 且 總和+1 整除10 返回 true
	case GUINumber[6] == '7' && (sum+1)%10 == 0:
		return true
	default:
		return false
	}
}

// IsNum 檢查字串是否全是數字
func IsNum(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
