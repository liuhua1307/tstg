package controllers

import (
	"log"
	"tangsong-esports/database"
	"tangsong-esports/models"
	"tangsong-esports/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// Login 用户登录
// @Summary 用户登录
// @Description 内部成员登录接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "登录信息"
// @Success 200 {object} models.Response{data=models.LoginResponse}
// @Router /api/v1/login [post]
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误")
		return
	}

	// 查找用户
	var member models.InternalMember
	result := database.DB.Where("account = ?", req.Account).First(&member)
	if result.Error != nil {
		utils.Error(c, "账号或密码错误")
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, member.PasswordHash) {
		utils.Error(c, "账号或密码错误")
		return
	}

	// 检查账户状态
	if member.Status != "正常" {
		utils.Error(c, "账户已被禁用")
		return
	}

	if !member.IsEnabled {
		utils.Error(c, "账户未启用")
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(member.MemberID, member.Account, member.UserRole)
	if err != nil {
		log.Printf("[Login] 生成访问令牌失败: %v", err)
		utils.Error(c, "生成令牌失败")
		return
	}

	// 生成刷新令牌
	refreshToken, err := utils.GenerateRefreshToken(member.MemberID, member.Account, member.UserRole)
	if err != nil {
		log.Printf("[Login] 生成刷新令牌失败: %v", err)
		utils.Error(c, "生成刷新令牌失败")
		return
	}

	// 记录登录日志
	loginLog := models.MemberLoginLogs{
		MemberID:    member.MemberID,
		LoginTime:   time.Now(),
		LoginIP:     c.ClientIP(),
		UserAgent:   c.GetHeader("User-Agent"),
		LoginStatus: "成功",
	}
	database.DB.Create(&loginLog)

	// 更新最后登录信息
	now := time.Now()
	member.LastLoginAt = &now
	member.LastLoginIP = c.ClientIP()
	database.DB.Save(&member)

	log.Printf("[Login] 用户登录成功: %s (ID: %d)", member.Account, member.MemberID)

	// 返回登录信息
	response := models.RefreshTokenResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(24 * 3600), // 24小时
		User:         &member,
	}

	utils.SuccessWithMessage(c, "登录成功", response)
}

// RefreshToken 刷新访问令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param refresh body models.RefreshTokenRequest true "刷新令牌信息"
// @Success 200 {object} models.Response{data=models.RefreshTokenResponse}
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Router /api/v1/refresh [post]
func RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[RefreshToken] 请求参数错误: %v", err)
		utils.Error(c, "请求参数错误")
		return
	}

	if req.RefreshToken == "" {
		utils.Unauthorized(c, "刷新令牌不能为空")
		return
	}

	log.Printf("[RefreshToken] 开始解析刷新令牌")

	// 解析刷新令牌
	claims, err := utils.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		log.Printf("[RefreshToken] 令牌解析失败: %v", err)
		utils.Unauthorized(c, "刷新令牌无效或已过期")
		return
	}

	log.Printf("[RefreshToken] 令牌解x析成功，用户ID: %d", claims.MemberID)

	// 查找用户
	var member models.InternalMember
	result := database.DB.Where("member_id = ?", claims.MemberID).First(&member)
	if result.Error != nil {
		log.Printf("[RefreshToken] 用户查询失败，用户ID: %d, 错误: %v", claims.MemberID, result.Error)
		utils.Unauthorized(c, "用户不存在")
		return
	}

	// 检查账户状态
	if member.Status != "正常" {
		log.Printf("[RefreshToken] 账户状态异常，用户ID: %d, 状态: %s", claims.MemberID, member.Status)
		utils.Unauthorized(c, "账户已被禁用")
		return
	}

	if !member.IsEnabled {
		log.Printf("[RefreshToken] 账户未启用，用户ID: %d", claims.MemberID)
		utils.Unauthorized(c, "账户未启用")
		return
	}

	// 生成新的访问令牌
	accessToken, err := utils.GenerateToken(member.MemberID, member.Account, member.UserRole)
	if err != nil {
		log.Printf("[RefreshToken] 生成访问令牌失败: %v", err)
		utils.Error(c, "生成访问令牌失败")
		return
	}

	// 生成新的刷新令牌
	refreshToken, err := utils.GenerateRefreshToken(member.MemberID, member.Account, member.UserRole)
	if err != nil {
		log.Printf("[RefreshToken] 生成刷新令牌失败: %v", err)
		utils.Error(c, "生成刷新令牌失败")
		return
	}

	log.Printf("[RefreshToken] 令牌刷新成功，用户ID: %d", claims.MemberID)

	// 返回新的令牌信息
	response := models.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(24 * 3600), // 24小时
		User:         &member,
	}

	utils.SuccessWithMessage(c, "令牌刷新成功", response)
}
