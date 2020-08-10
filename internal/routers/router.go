package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-server/global"
	"github.com/go-programming-tour-book/blog-server/internal/middleware"
	"github.com/go-programming-tour-book/blog-server/internal/routers/api"
	v1 "github.com/go-programming-tour-book/blog-server/internal/routers/api/v1"
	"github.com/go-programming-tour-book/blog-server/pkg/limiter"
	_ "github.com/go-programming-tour-book/blog-service/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"time"
)

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(
	limiter.LimiterBucketRule{
		Key:          "/auth",
		FillInterval: time.Second,
		Capacity:     10,
		Quantum:      10,
	},
)

func NewRouter() *gin.Engine {
	r := gin.New()

	if global.SeverSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog())
		r.Use(middleware.Recovery())
	}

	r.Use(middleware.RateLimiter(methodLimiters))
	r.Use(middleware.ContextTimeout(global.SeverSetting.ContextTimeout))
	r.Use(middleware.Translations())
	r.Use(middleware.Tracing())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/auth", api.GetAuth)

	tag := v1.NewTag()

	article := v1.NewArticle()

	upload := api.NewUpload()

	r.POST("/upload/file", upload.UploadFile)
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))

	apiV1 := r.Group("/api/v1")

	apiV1.Use(middleware.JWT())
	{
		apiV1.POST("/tags", tag.Create)
		apiV1.DELETE("/tags/:id", tag.Delete)
		apiV1.PUT("/tags/:id", tag.Update)
		apiV1.PATCH("/tags/:id/state", tag.Update)
		apiV1.GET("/tags/:id", tag.Get)
		apiV1.GET("/tags", tag.List)

		apiV1.POST("/articles", article.Create)
		apiV1.DELETE("/articles/:id", article.Delete)
		apiV1.PUT("/articles/:id", article.Update)
		apiV1.PATCH("/articles/:id/state", article.Update)
		apiV1.GET("/articles/:id", article.Get)
		apiV1.GET("/articles", article.List)
	}

	//article := v1.NewArticle()
	//tag := v1.NewTag()
	//apiV1 := r.Group("/api/v1")
	return r
}
