package router

import (
	"tangsong-esports/controllers"
	"tangsong-esports/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	// 中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS 中间件
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 公开路由（无需认证）
		public := api.Group("")
		{
			// 内部员工认证
			public.POST("/login", controllers.Login)
			public.POST("/refresh", controllers.RefreshToken)

			// 客户认证
			public.POST("/customer/register", controllers.CustomerRegister)
			public.POST("/customer/login", controllers.CustomerLogin)
			public.POST("/customer/refresh", controllers.CustomerRefreshToken)
			public.POST("/customer/set-password", controllers.CustomerSetPassword)
		}

		// 受保护路由（需要认证）
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// 内部成员管理
			members := protected.Group("/members")
			{
				members.GET("", controllers.GetMembers)
				members.POST("", controllers.CreateMember)
				members.PUT("/:id", controllers.UpdateMember)
				members.DELETE("/:id", controllers.DeleteMember)
				members.GET("/:id", controllers.GetMemberByID)
			}

			// 客户管理
			customers := protected.Group("/customers")
			{
				customers.GET("", controllers.GetCustomers)
				customers.POST("", controllers.CreateCustomer)
				customers.PUT("/:id", controllers.UpdateCustomer)
				customers.DELETE("/:id", controllers.DeleteCustomer)
				customers.GET("/:id", controllers.GetCustomerByID)
				customers.POST("/:id/recharge", controllers.RechargeCustomer)
				customers.GET("/:id/recharge-history", controllers.GetCustomerRechargeHistory)
				customers.POST("/reset-password", controllers.AdminResetCustomerPassword)
			}

			// 客户个人信息
			protected.GET("/customer/profile", controllers.GetCustomerProfile)
			protected.PUT("/customer/profile", controllers.UpdateCustomerProfile)
			protected.GET("/customer/balance", controllers.GetCustomerBalance)
			protected.POST("/customer/reset-password", controllers.CustomerResetPassword)
			// 客户查看自己的订单
			protected.GET("/customer/orders", controllers.GetCustomerOrders)

			// 订单类别管理
			categories := protected.Group("/order-categories")
			{
				categories.GET("", controllers.GetOrderCategories)
				categories.POST("", controllers.CreateOrderCategory)
				categories.PUT("/:id", controllers.UpdateOrderCategory)
				categories.DELETE("/:id", controllers.DeleteOrderCategory)
			}

			// 陪玩报单管理
			orders := protected.Group("/orders")
			{
				orders.GET("", controllers.GetOrders)
				orders.POST("", controllers.CreateOrder)
				orders.PUT("/:id", controllers.UpdateOrder)
				orders.GET("/:id", controllers.GetOrderByID)
				orders.PUT("/:id/status", controllers.UpdateOrderStatus)
				orders.POST("/:id/images", controllers.UploadOrderImages)
				orders.GET("/stats", controllers.GetOrderStats)
			}

			// 订单审批管理
			approval := protected.Group("/order-approval")
			{
				approval.GET("/pending", controllers.GetPendingOrders)
				approval.GET("", controllers.GetApprovalOrders)
				approval.POST("/:id/approve", controllers.ApproveOrder)
				approval.POST("/:id/reject", controllers.RejectOrder)
				approval.POST("/batch", controllers.BatchApproval)
				approval.PATCH("/:id/status", controllers.UpdateOrderStatusV2)
				approval.GET("/statistics", controllers.GetStatistics)
				approval.GET("/history", controllers.GetOperationHistory)
				approval.POST("/export", controllers.BatchExport)
			}

			// 系统配置
			configs := protected.Group("/configs")
			{
				configs.GET("", controllers.GetConfigs)
				configs.PUT("/:key", controllers.UpdateConfig)
			}

			// 操作日志
			logs := protected.Group("/logs")
			{
				logs.GET("", controllers.GetOperationLogs)
			}
		}
	}

	return r
}
