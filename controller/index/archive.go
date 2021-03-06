package index

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/model"
	"github.com/zwh8800/md-blog-gen/render"
	"github.com/zwh8800/md-blog-gen/service"
	"github.com/zwh8800/md-blog-gen/util"
)

func Archive(c *gin.Context) {
	monthList, noteListMap, err := service.NoteGroupByMonth()
	if err != nil {
		glog.Error(err)
		ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		return
	}

	c.Render(http.StatusOK, render.NewRender("archive.html", gin.H{
		"monthList":   monthList,
		"noteListMap": noteListMap,
		"site":        conf.Conf.Site,
		"social":      conf.Conf.Social,
		"prod":        conf.Conf.Env.Prod,
		"haha":        util.HahaGenarate(),
	}))
}

func ArchiveMonth(c *gin.Context) {
	monthStr := c.Param("month")
	timeMonth, err := time.Parse("2006-1", monthStr)
	if err != nil {
		glog.Error(err)
		ErrorHandler(c, http.StatusNotFound, errors.New("Not Found"))
		return
	}
	month := &model.YearMonth{
		int64(timeMonth.Year()),
		int64(timeMonth.Month()),
	}

	noteList, err := service.NotesByMonth(month)
	if err != nil {
		glog.Error(err)
		ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		return
	}
	if len(noteList) == 0 {
		ErrorHandler(c, http.StatusNotFound, errors.New("Not Found"))
		return
	}

	prev, next, err := service.PrevNextMonth(month)
	if err != nil {
		glog.Error(err)
		ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		return
	}

	c.Render(http.StatusOK, render.NewRender("archive_month.html", gin.H{
		"prevMonth": prev,
		"nextMonth": next,
		"month":     month,
		"noteList":  noteList,
		"site":      conf.Conf.Site,
		"social":    conf.Conf.Social,
		"prod":      conf.Conf.Env.Prod,
		"haha":      util.HahaGenarate(),
	}))
}
