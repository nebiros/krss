package middleware

import (
	"html/template"
	"io"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gosimple/slug"

	"github.com/dustin/go-humanize"

	"github.com/pkg/errors"

	"github.com/labstack/echo/v4"
	"github.com/oxtoacart/bpool"
)

const (
	defaultTemplate = `{{define "main"}} {{template "base" .}} {{end}}`
)

type fileData struct {
	absolutePath string
	fileName     string
	basename     string
	key          string
}

type IncludeData struct {
	Title   string
	Flashes []interface{}
	Data    interface{}
}

type TemplateRendererConfig struct {
	LayoutPattern  string
	IncludePattern string
	Templates      map[string]*template.Template
	BufferPool     *bpool.BufferPool
}

var (
	DefaultTemplateRendererConfig = TemplateRendererConfig{
		LayoutPattern:  "../../web/template/layout/*.gohtml",
		IncludePattern: "../../web/template/include/**/*.gohtml",
		Templates:      make(map[string]*template.Template),
		BufferPool:     bpool.NewBufferPool(64),
	}
)

type TemplateRenderer struct {
	layoutPattern  string
	includePattern string
	includeFiles   map[string]fileData
	templates      map[string]*template.Template
	bufpool        *bpool.BufferPool
}

func NewTemplateRenderer() (*TemplateRenderer, error) {
	return NewTemplateRendererWithConfig(DefaultTemplateRendererConfig)
}

func NewTemplateRendererWithConfig(config TemplateRendererConfig) (*TemplateRenderer, error) {
	t := &TemplateRenderer{
		layoutPattern:  config.LayoutPattern,
		includePattern: config.IncludePattern,
		includeFiles:   make(map[string]fileData),
		templates:      config.Templates,
		bufpool:        config.BufferPool,
	}

	if err := t.loadTemplates(); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *TemplateRenderer) loadTemplates() error {
	layoutFiles, err := filepath.Glob(t.layoutPattern)
	if err != nil {
		return err
	}

	includeFiles, err := filepath.Glob(t.includePattern)
	if err != nil {
		return err
	}

	mainTemplate, err := template.New("main").Parse(defaultTemplate)
	if err != nil {
		return err
	}

	for _, file := range includeFiles {
		fd := t.makeFileData(file)

		t.includeFiles[fd.key] = fd

		if strings.HasPrefix(fd.basename, "_") {
			continue
		}

		files := append(layoutFiles, file)

		clone, err := mainTemplate.Clone()
		if err != nil {
			return errors.Wrap(err, "template clone failed")
		}

		t.templates[fd.key] = clone.Funcs(template.FuncMap{
			"include":       t.includeControl,
			"timeFormat":    t.timeFormatControl,
			"timeHumanized": t.timeHumanizedControl,
			"htmlSafe":      t.htmlSafeControl,
			"slug":          t.slugControl,
		})
		t.templates[fd.key] = template.Must(t.templates[fd.key].ParseFiles(files...))
	}

	return nil
}

func (t *TemplateRenderer) makeFileData(absolutePath string) fileData {
	basename := filepath.Base(absolutePath)
	fileName := strings.TrimSuffix(basename, path.Ext(basename))
	key := filepath.Base(filepath.Dir(absolutePath)) + "/" + fileName

	return fileData{
		absolutePath: absolutePath,
		fileName:     fileName,
		basename:     basename,
		key:          key,
	}
}

func (t *TemplateRenderer) includeControl(key string, contextList ...interface{}) (interface{}, error) {
	fd, exists := t.includeFiles[key]
	if !exists {
		return nil, errors.Errorf("partial '%s' not found", key)
	}

	tpl := template.Must(template.New(fd.key).ParseFiles(fd.absolutePath))

	buf := t.bufpool.Get()
	defer t.bufpool.Put(buf)

	var data IncludeData

	if len(contextList) == 0 {
		data = IncludeData{}
	} else {
		data = contextList[0].(IncludeData)
	}

	if err := tpl.Execute(buf, data); err != nil {
		return nil, errors.Wrapf(err, "partial '%s' execution failed", key)
	}

	return template.HTML(buf.String()), nil
}

func (t *TemplateRenderer) timeFormatControl(ti time.Time, layout string) string {
	return ti.Format(layout)
}

func (t *TemplateRenderer) timeHumanizedControl(ti time.Time) string {
	return humanize.Time(ti)
}

func (t *TemplateRenderer) htmlSafeControl(html string) template.HTML {
	return template.HTML(html)
}

func (t *TemplateRenderer) slugControl(s string) string {
	return slug.Make(s)
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	tpl, exists := t.templates[name]
	if !exists {
		return errors.Errorf("template '%s' not found", name)
	}

	buf := t.bufpool.Get()
	defer t.bufpool.Put(buf)

	if err := tpl.Execute(buf, data.(IncludeData)); err != nil {
		return errors.Wrapf(err, "template '%s' execution failed", name)
	}

	if _, err := buf.WriteTo(w); err != nil {
		return errors.Wrapf(err, "template '%s' write failed", name)
	}

	return nil
}
