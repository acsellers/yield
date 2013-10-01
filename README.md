yield
=====

This library allows you to use yields and layouts similar to the Rails implementation for the current Revel template implementation.

Instruction for Use
-------------------

In your app.conf, add a line like
module.yield=github.com/acsellers/yield

Then instead of starting your controllers from \*revel.Controller, 
you can import "github.com/acsellers/yield/app/controllers" and then
use the struct yield.Controller to embed into your controllers. Note:
the module in that import path is named yield not controllers, and 
that is why you embed yield.Controller not controllers.Controller.

The booking sample from revel was ported to use the basic yield
mechanism and is available in the samples directory.

Documentation is at [godoc.org/github.com/acsellers/yield/app/controllers](godoc.org/github.com/acsellers/yield/app/controllers).

Bugs
----

Please file an issue with details and I will get on it. If you would 
prefer to submit a pull request that fixes it, that would be acceptable
as well.

Future Features
---------------

Yield has it's single feature, and I'm fine with that. Now that it
functions in the manner I would like, I'm moving onto its bit brother
unitemplate. Hopefully unitemplate will support 10+ html template 
formats, a few xml and json formats, plus asset pipelines.
