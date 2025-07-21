package main
import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Claims struct {
	MemberID  uint   `json:"member_id"`
	Account   string `json:"account"`
	UserRole  string `json:"user_role"`
	TokenType string `json:"token_type"`
	Exp       int64  `json:"exp"`
	Iat       int64  `json:"iat"`
}
func createSignature(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func parseClaimsManually(jsonStr string) (*Claims, error) {
	claims := &Claims{}
	
	// 去掉花括号
	jsonStr = strings.Trim(jsonStr, "{}")
	
	// 按引号分割并解析
	parts := strings.Split(jsonStr, `"`)
	for i := 0; i < len(parts)-1; i += 4 {
		if i+3 >= len(parts) {
			break
		}
		
		key := parts[i+1]
		value := parts[i+3]
		
		switch key {
		case "member_id":
			if val, err := strconv.ParseUint(value, 10, 32); err == nil {
				claims.MemberID = uint(val)
			}
		case "account":
			claims.Account = value
		case "user_role":
			claims.UserRole = value
		case "token_type":
			claims.TokenType = value
		case "exp":
			if val, err := strconv.ParseInt(value, 10, 64); err == nil {
				claims.Exp = val
			}
		case "iat":
			if val, err := strconv.ParseInt(value, 10, 64); err == nil {
				claims.Iat = val
			}
		}
	}
	
	return claims, nil
}

func debugToken(tokenString string) {
	secret := "tangsong-esports-secret-key"
	
	fmt.Printf("调试令牌: %s\n", tokenString)
	fmt.Printf("令牌长度: %d\n", len(tokenString))
	
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		fmt.Printf("❌ 令牌格式错误，部分数量: %d\n", len(parts))
		return
	}
	
	fmt.Printf("✅ 令牌分段正确，数量: %d\n", len(parts))
	fmt.Printf("Header 长度: %d, Payload 长度: %d, Signature 长度: %d\n", 
		len(parts[0]), len(parts[1]), len(parts[2]))
	
	// 解码 payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Printf("❌ Payload 解码失败: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Payload 解码成功\n")
	fmt.Printf("Payload 内容: %s\n", string(payloadBytes))
	
	// 尝试标准 JSON 解析
	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		fmt.Printf("❌ 标准 JSON 解析失败: %v\n", err)
		
		// 尝试手动解析
		if manualClaims, manualErr := parseClaimsManually(string(payloadBytes)); manualErr == nil {
			claims = *manualClaims
			fmt.Printf("✅ 手动解析成功\n")
		} else {
			fmt.Printf("❌ 手动解析也失败: %v\n", manualErr)
			return
		}
	} else {
		fmt.Printf("✅ 标准 JSON 解析成功\n")
	}
	
	fmt.Printf("解析结果: %+v\n", claims)
	
	// 验证签名
	message := parts[0] + "." + parts[1]
	expectedSignature := createSignature(message, secret)
	
	fmt.Printf("期望签名: %s\n", expectedSignature)
	fmt.Printf("实际签名: %s\n", parts[2])
	
	if expectedSignature == parts[2] {
		fmt.Printf("✅ 签名验证成功\n")
	} else {
		fmt.Printf("❌ 签名验证失败\n")
	}
	
	// 检查时间
	currentTime := time.Now().Unix()
	fmt.Printf("当前时间: %d\n", currentTime)
	fmt.Printf("令牌过期时间: %d\n", claims.Exp)
	
	if currentTime <= claims.Exp {
		fmt.Printf("✅ 令牌未过期\n")
	} else {
		fmt.Printf("❌ 令牌已过期\n")
	}
	// 检查令牌类型
	if claims.TokenType == "refresh" {
		fmt.Printf("✅ 令牌类型正确: refresh\n")
	} else {
		fmt.Printf("❌ 令牌类型错误: %s\n", claims.TokenType)
	}
}

func main() {
	// 测试从登录响应中获取的 refresh_token
	refreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtZW1iZXJfaWQiOjEsImFjY291bnQiOiJhZG1pbiIsInVzZXJfcm9sZSI6Iui2hee6p-euoeeQhuWRmCIsInRva2VuX3R5cGUiOiJyZWZyZXNoIiwiZXhwIjoxNzUzMzI3ODc2LCJpYXQiOjE3NTI3MjMwNzZ9.dg6iI8pOR3qL-GF8V3tgmz2pUebG0FwTMkogFCBZtYQ"
	
	fmt.Println("=== 调试 Refresh Token ===")
	debugToken(refreshToken)
}
