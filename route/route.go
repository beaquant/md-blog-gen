package route

import (
	"errors"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/zwh8800/md-blog-gen/index"
	"github.com/zwh8800/md-blog-gen/rss"
	"github.com/zwh8800/md-blog-gen/sitemap"
	"github.com/zwh8800/md-blog-gen/util"
)

func Route(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		index.ErrorHandler(c, 404, errors.New("404 Not Found"))
	})

	indexGroup := r.Group("/")
	{
		indexGroup.GET("/", index.Index)
		indexGroup.GET(path.Join(util.GetPageBase(), ":page"), index.Index)
		indexGroup.GET(util.GetTagBase(), index.AllTag)
		indexGroup.GET(path.Join(util.GetTagBase(), ":tag"), index.Tag)
		indexGroup.GET(path.Join(util.GetNoteBase(), ":id"), index.Note)
		indexGroup.GET(util.GetArchiveBase(), index.Archive)
		indexGroup.GET(path.Join(util.GetArchiveBase(), ":month"), index.ArchiveMonth)
	}
	rssGroup := r.Group(util.GetRssBase())
	{
		rssGroup.GET("", rss.Rss)
	}
	sitemapGroup := r.Group("/sitemap.xml")
	{
		sitemapGroup.GET("", sitemap.SiteMap)
	}

	r.Static("/static", "./static")
}
