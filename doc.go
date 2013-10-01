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
