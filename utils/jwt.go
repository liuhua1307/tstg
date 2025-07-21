package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"tangsong-esports/config"
	"time"
)

type Claims struct {
	MemberID  uint   `json:"member_id"`
	Account   string `json:"account"`
	UserRole  string `json:"user_role"`
	TokenType string `json:"token_type"` // access 或 refresh
	Exp       int64  `json:"exp"`
	Iat       int64  `json:"iat"`
}

// GenerateToken 生成JWT令牌 (简化版本)
func GenerateToken(memberID uint, account, userRole string) (string, error) {
	return generateToken(memberID, account, userRole, "access", 24*time.Hour)
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(memberID uint, account, userRole string) (string, error) {
	return generateToken(memberID, account, userRole, "refresh", 7*24*time.Hour) // 7天
}

// generateToken 生成指定类型的令牌 - 使用字符串拼接避免 JSON 序列化问题
func generateToken(memberID uint, account, userRole, tokenType string, duration time.Duration) (string, error) {
	// JWT Header - 手动构建 JSON 字符串
	headerJSON := `{"alg":"HS256","typ":"JWT"}`
	headerEncoded := base64.RawURLEncoding.EncodeToString([]byte(headerJSON))

	// JWT Payload - 手动构建 JSON 字符串
	now := time.Now()
	exp := now.Add(duration).Unix()
	iat := now.Unix()
	// 转义用户角色中的特殊字符
	escapedUserRole := strings.ReplaceAll(userRole, `"`, `\"`)
	payloadJSON := fmt.Sprintf(`{"member_id":%d,"account":"%s","user_role":"%s","token_type":"%s","exp":%d,"iat":%d}`,
		memberID, account, escapedUserRole, tokenType, exp, iat)
	payloadEncoded := base64.RawURLEncoding.EncodeToString([]byte(payloadJSON))

	// 创建签名
	message := headerEncoded + "." + payloadEncoded
	signature := createSignature(message, config.AppConfig.JWT.Secret)

	return message + "." + signature, nil
}

// ParseToken 解析JWT令牌 (简化版本)
func ParseToken(tokenString string) (*Claims, error) {
	claims, err := parseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 检查令牌类型
	if claims.TokenType != "access" {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}

// ParseRefreshToken 解析刷新令牌
func ParseRefreshToken(tokenString string) (*Claims, error) {
	claims, err := parseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 检查令牌类型
	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("invalid refresh token")
	}

	return claims, nil
}

// parseToken 解析令牌的通用函数
func parseToken(tokenString string) (*Claims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// 验证签名
	message := parts[0] + "." + parts[1]
	expectedSignature := createSignature(message, config.AppConfig.JWT.Secret)
	if parts[2] != expectedSignature {
		log.Printf("[parseToken] 签名验证失败 - 期望: %s, 实际: %s", expectedSignature, parts[2])
		return nil, fmt.Errorf("invalid signature")
	}

	// 解码 Payload
	payloadBytes, err := base64URLDecode(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		// 尝试手动解析 JSON（如果标准解析失败）
		if manualClaims, manualErr := parseClaimsManually(string(payloadBytes)); manualErr == nil {
			claims = *manualClaims
		} else {
			return nil, fmt.Errorf("failed to unmarshal claims: %v", err)
		}
	}

	// 检查过期时间
	currentTime := time.Now().Unix()
	if currentTime > claims.Exp {
		log.Printf("[parseToken] 令牌已过期 - 当前: %d, 过期: %d", currentTime, claims.Exp)
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}

// parseClaimsManually 手动解析 Claims（用于处理格式不正确的 JSON）
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

// base64URLDecode 兼容标准 Base64URL 解码
func base64URLDecode(data string) ([]byte, error) {
	// 首先尝试 RawURLEncoding
	if decoded, err := base64.RawURLEncoding.DecodeString(data); err == nil {
		return decoded, nil
	}
	// 如果失败，尝试标准 URLEncoding
	if decoded, err := base64.URLEncoding.DecodeString(data); err == nil {
		return decoded, nil
	}
	// 如果还是失败，尝试添加填充后再解码
	switch len(data) % 4 {
	case 2:
		data += "=="
	case 3:
		data += "="
	}
	return base64.URLEncoding.DecodeString(data)
}

// base64URLEncode 标准 Base64URL 编码
func base64URLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// createSignature 创建HMAC-SHA256签名
func createSignature(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
