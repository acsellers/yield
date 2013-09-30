// Yield provides a custom Controller for revel apps
// that allows you to use layouts and output varying
// templates from templates.
package yield

import (
	"bytes"
	"fmt"
	"github.com/robfig/revel"
	htmlTmpl "html/template"
	"path/filepath"
	"runtime"
)

var (
	LayoutPath      = "app/layouts"
	DefaultLayout   = make(map[string]string)
	layoutTemplates *revel.TemplateLoader
)

func init() {
	revel.TemplateFuncs["could_yield"] = func(name string, renderArgs map[string]interface{}) bool {
		if items, found := renderArgs["ContentForItems"]; found {
			if renderTmpl, ok := items.(map[string]revel.Template); ok {
				_, found := renderTmpl[name]
				return found
			} else {
				return false
			}
		} else {
			return false
		}
	}

	revel.TemplateFuncs["yield"] = func(args ...interface{}) (htmlTmpl.HTML, error) {
		var renderArgs map[string]interface{}
		var target string

		switch len(args) {
		case 1:
			if r_arg, ok := args[0].(map[string]interface{}); ok {
				renderArgs = r_arg
			} else {
				return "", fmt.Errorf("Must pass dot into yield")
			}
		case 2:
			if t_arg, ok := args[0].(string); ok {
				target = t_arg
			} else {
				return "", fmt.Errorf("Named yields require the name as the first argument")
			}
			if r_arg, ok := args[1].(map[string]interface{}); ok {
				renderArgs = r_arg
			} else {
				return "", fmt.Errorf("Named yields require the dot as the second argument")
			}
		default:
			return "", fmt.Errorf("Yield: Argument Length Error")
		}

		if items, found := renderArgs["ContentForItems"]; found {
			if renderTmpl, ok := items.(map[string]revel.Template); ok {
				if tmpl, found := renderTmpl[target]; found {
					var b bytes.Buffer
					err := tmpl.Render(&b, renderArgs)
					if err != nil {
						return "", err
					}
					return htmlTmpl.HTML(b.String()), nil
				}
				return "", nil
			} else {
				return "", fmt.Errorf("Yield: ContentForItems was overwritten")
			}
		} else {
			return "", fmt.Errorf("Yield requires the base RenderArgs")
		}
	}
}

type LayoutController struct {
	*revel.Controller
	RenderTmpl map[string]revel.Template
	LayoutPath string
	noLayout   bool
}

func (lc *LayoutController) Layout(s string) {
	if s == "" {
		lc.noLayout = true
	} else {
		lc.LayoutPath = s
	}
}

func (lc *LayoutController) Render(extraRenderArgs ...interface{}) revel.Result {
	// Get the calling function name.
	_, _, line, ok := runtime.Caller(1)
	if !ok {
		revel.ERROR.Println("Failed to get Caller information")
	}

	// Get the extra RenderArgs passed in.
	if renderArgNames, ok := lc.MethodType.RenderArgNames[line]; ok {
		if len(renderArgNames) == len(extraRenderArgs) {
			for i, extraRenderArg := range extraRenderArgs {
				lc.RenderArgs[renderArgNames[i]] = extraRenderArg
			}
		} else {
			revel.ERROR.Println(len(renderArgNames), "RenderArg names found for",
				len(extraRenderArgs), "extra RenderArgs")
		}
	} else {
		revel.ERROR.Println("No RenderArg names found for Render call on line", line,
			"(Method", lc.MethodType.Name, ")")
	}
	if lc.noLayout || DefaultLayout[lc.Request.Format] == "" {
		return lc.RenderTemplate(lc.Name + "/" + lc.MethodType.Name + "." + lc.Request.Format)
	} else {
		if lc.LayoutPath == "" {
			lc.LayoutPath = DefaultLayout[lc.Request.Format]
		}
		return lc.RenderTemplateWithLayout(lc.Name + "/" + lc.MethodType.Name + "." + lc.Request.Format)
	}
}

func (lc *LayoutController) RenderTemplateWithLayout(templatePath string) revel.Result {
	if layoutTemplates == nil {
		err := loadLayouts()
		if err != nil {
			return lc.RenderError(err)
		}
	}

	// Get the Template.
	template, err := revel.MainTemplateLoader.Template(templatePath)
	if err != nil {
		return lc.RenderError(err)
	}
	layout, err := layoutTemplates.Template(lc.LayoutPath)
	if err != nil {
		return lc.RenderError(err)
	}

	return &RenderLayoutTemplateResult{
		Template:   template,
		Layout:     layout,
		RenderArgs: lc.RenderArgs,
	}
}

func loadLayouts() *revel.Error {
	layoutTemplates := revel.NewTemplateLoader([]string{filepath.Join(revel.BasePath, LayoutPath)})
	return layoutTemplates.Refresh()
}

// Set a template from your main revel Template library to be rendered into
// a named yield.
func (lc *LayoutController) ContentFor(yieldName, templateName string) error {
	template, err := revel.MainTemplateLoader.Template(templateName)
	if err != nil {
		template, err = revel.MainTemplateLoader.Template(lc.Name + "/" + templateName)
		if err != nil {
			return err
		}
	}

	lc.RenderTmpl[yieldName] = template

	return nil
}
