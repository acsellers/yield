/*
This library allows you to use yields and layouts similar to the Rails implementation for the current Revel template implementation.

In your app.conf, add a line like
module.yield=github.com/acsellers/yield

Then instead of starting your controllers from *revel.Controller,
you can import "github.com/acsellers/yield/app/controllers" and then
use the struct yield.Controller to embed into your controllers.

Note: the module in that import path is named yield not controllers, and
that is why you embed yield.Controller not controllers.Controller.

The booking sample from revel was ported to use the basic yield
mechanism and is available in the samples directory.
*/
package yield

import (
	"bytes"
	"fmt"
	"github.com/robfig/revel"
	htmlTmpl "html/template"
	"path/filepath"
	"runtime"
)

/*
To set your directory for loading layouts from, set the LayoutPath variable to the location
relative to the base of your revel directory. You can only set one directory at the moment.

To set a default layout, take the format you wish that layout to apply for, i.e. "html", then set
that string to the name of the layout you want to render.
*/
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

/*
You can embed this Controller into your controllers instead of *revel.Controller,
note that unlike revel.Controller, you do not need to embed a pointer to this
controller.
*/
type Controller struct {
	*revel.Controller
	RenderTmpl map[string]revel.Template
	LayoutPath string
	noLayout   bool
}

/*
Set the layout to be rendered for the current action. Setting the layout
to empty string will cause no layout to be rendered. No layout will be rendered
if you did not set a Default layout for the current request format and you
do not call this function to set a specific layout. You do not have to include
the format for the template, that will be added, but adding it the format would
not cause a problem.
*/
func (lc *Controller) Layout(s string) {
	if s == "" {
		lc.noLayout = true
	} else {
		lc.LayoutPath = s
	}
}

/*
The same kind of function as calling Render on revel.Controller, except that
this will pick up the Layout you specified and render that as well.
*/
func (lc *Controller) Render(extraRenderArgs ...interface{}) revel.Result {
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

/*
If you needed to use revel's RenderTemplate, this is similar, except it uses
the Layout specified on the Controller. If you do not wish for a Layout to be
rendered, you should use RenderTemplate which is available. This call expects
an actual layout to be rendered, failing to provide one is an error.
*/
func (lc *Controller) RenderTemplateWithLayout(templatePath string) revel.Result {
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
		layout, err = layoutTemplates.Template(lc.LayoutPath + "." + lc.Request.Format)
		if err != nil {
			return lc.RenderError(err)
		}
	}

	return &RenderLayoutTemplateResult{
		Template:   template,
		Layout:     layout,
		RenderArgs: lc.RenderArgs,
		RenderTmpl: map[string]revel.Template{},
	}
}

func loadLayouts() *revel.Error {
	layoutTemplates = revel.NewTemplateLoader([]string{filepath.Join(revel.BasePath, LayoutPath)})
	return layoutTemplates.Refresh()
}

// Set a template from your main revel Template library to be rendered into
// a named yield.
func (lc *Controller) ContentFor(yieldName, templateName string) error {
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
