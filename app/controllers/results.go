package yield

import (
	"bytes"
	"fmt"
	"github.com/revel/revel"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

/*
This is the Result object passed to revel when you render
a view with a Layout set. RenderTmpl is the Templates that
were set with ContentFor. Layout is the set Layout. Otherwise
it is the same as revel's RenderTemplateResult. This actually
doesn't require a Layout to be set, not that its used with that
functionality.
*/
type RenderLayoutTemplateResult struct {
	Template   revel.Template
	Layout     revel.Template
	RenderArgs map[string]interface{}
	RenderTmpl map[string]revel.Template
}

// Render the Templates into the Response, handles errors and panics using the
// same mechanisms of revel.
func (r *RenderLayoutTemplateResult) Apply(req *revel.Request, resp *revel.Response) {
	// Handle panics when rendering templates.
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Println(err)
			revel.PlaintextErrorResult{fmt.Errorf("Template Execution Panic in %s:\n%s",
				r.Template.Name(), err)}.Apply(req, resp)
		}
	}()

	chunked := revel.Config.BoolDefault("results.chunked", false)
	r.RenderTmpl[""] = r.Template
	r.RenderArgs["ContentForItems"] = r.RenderTmpl

	// If it's a HEAD request, throw away the bytes.
	out := io.Writer(resp.Out)
	if req.Method == "HEAD" {
		out = ioutil.Discard
	}

	// In a prod mode, write the status, render, and hope for the best.
	// (In a dev mode, always render to a temporary buffer first to avoid having
	// error pages distorted by HTML already written)
	if chunked && !revel.DevMode {
		resp.WriteHeader(http.StatusOK, "text/html")
		if r.Layout == nil {
			r.render(req, resp, out)
		} else {
			r.renderWithLayout(req, resp, out)
		}
		return
	}

	// Render the template into a temporary buffer, to see if there was an error
	// rendering the template.  If not, then copy it into the response buffer.
	// Otherwise, template render errors may result in unpredictable HTML (and
	// would carry a 200 status code)
	var b bytes.Buffer
	if r.Layout == nil {
		r.render(req, resp, &b)
	} else {
		r.renderWithLayout(req, resp, &b)
	}

	if !chunked {
		resp.Out.Header().Set("Content-Length", strconv.Itoa(b.Len()))
	}
	resp.WriteHeader(http.StatusOK, "text/html")
	b.WriteTo(out)
}

func (r *RenderLayoutTemplateResult) render(req *revel.Request, resp *revel.Response, wr io.Writer) {
	err := r.Template.Render(wr, r.RenderArgs)
	if err == nil {
		return
	}
	r.renderError(req, resp, err)
}

func (r *RenderLayoutTemplateResult) renderWithLayout(req *revel.Request, resp *revel.Response, wr io.Writer) {
	err := r.Layout.Render(wr, r.RenderArgs)
	if err == nil {
		return
	}
	r.renderError(req, resp, err)
}

func (r *RenderLayoutTemplateResult) renderError(req *revel.Request, resp *revel.Response, err error) {
	var templateContent []string
	templateName, line, description := parseTemplateError(err)
	if templateName == "" {
		templateName = r.Layout.Name()
		templateContent = r.Layout.Content()
	} else {
		if tmpl, err := revel.MainTemplateLoader.Template(templateName); err == nil {
			templateContent = tmpl.Content()
		}
	}
	compileError := &revel.Error{
		Title:       "Layout Execution Error",
		Path:        templateName,
		Description: description,
		Line:        line,
		SourceLines: templateContent,
	}
	resp.Status = 500
	revel.ERROR.Printf("Template Execution Error (in %s): %s", templateName, description)
	revel.ErrorResult{r.RenderArgs, compileError}.Apply(req, resp)
}

// Parse the line, and description from an error message like:
// html/template:Application/Register.html:36: no such template "footer.html"
func parseTemplateError(err error) (templateName string, line int, description string) {
	description = err.Error()
	i := regexp.MustCompile(`:\d+:`).FindStringIndex(description)
	if i != nil {
		line, err = strconv.Atoi(description[i[0]+1 : i[1]-1])
		if err != nil {
			revel.ERROR.Println("Failed to parse line number from error message:", err)
		}
		templateName = description[:i[0]]
		if colon := strings.Index(templateName, ":"); colon != -1 {
			templateName = templateName[colon+1:]
		}
		templateName = strings.TrimSpace(templateName)
		description = description[i[1]+1:]
	}
	return templateName, line, description
}
