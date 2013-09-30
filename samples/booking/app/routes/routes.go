// GENERATED CODE - DO NOT EDIT
package routes

import "github.com/robfig/revel"


type tGorpController struct {}
var GorpController tGorpController


func (_ tGorpController) Begin(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("GorpController.Begin", args).Url
}

func (_ tGorpController) Commit(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("GorpController.Commit", args).Url
}

func (_ tGorpController) Rollback(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("GorpController.Rollback", args).Url
}


type tStatic struct {}
var Static tStatic


func (_ tStatic) Serve(
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.Serve", args).Url
}

func (_ tStatic) ServeModule(
		moduleName string,
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "moduleName", moduleName)
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.ServeModule", args).Url
}


type tJobs struct {}
var Jobs tJobs


func (_ tJobs) Status(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Jobs.Status", args).Url
}


type tTestRunner struct {}
var TestRunner tTestRunner


func (_ tTestRunner) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.Index", args).Url
}

func (_ tTestRunner) Run(
		suite string,
		test string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "suite", suite)
	revel.Unbind(args, "test", test)
	return revel.MainRouter.Reverse("TestRunner.Run", args).Url
}

func (_ tTestRunner) List(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.List", args).Url
}


type tApplication struct {}
var Application tApplication


func (_ tApplication) AddUser(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Application.AddUser", args).Url
}

func (_ tApplication) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Application.Index", args).Url
}

func (_ tApplication) Register(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Application.Register", args).Url
}

func (_ tApplication) SaveUser(
		user interface{},
		verifyPassword string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "user", user)
	revel.Unbind(args, "verifyPassword", verifyPassword)
	return revel.MainRouter.Reverse("Application.SaveUser", args).Url
}

func (_ tApplication) Login(
		username string,
		password string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	revel.Unbind(args, "password", password)
	return revel.MainRouter.Reverse("Application.Login", args).Url
}

func (_ tApplication) Logout(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Application.Logout", args).Url
}


type tHotels struct {}
var Hotels tHotels


func (_ tHotels) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Hotels.Index", args).Url
}

func (_ tHotels) List(
		search string,
		size int,
		page int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "search", search)
	revel.Unbind(args, "size", size)
	revel.Unbind(args, "page", page)
	return revel.MainRouter.Reverse("Hotels.List", args).Url
}

func (_ tHotels) Show(
		id int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("Hotels.Show", args).Url
}

func (_ tHotels) Settings(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Hotels.Settings", args).Url
}

func (_ tHotels) SaveSettings(
		password string,
		verifyPassword string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "password", password)
	revel.Unbind(args, "verifyPassword", verifyPassword)
	return revel.MainRouter.Reverse("Hotels.SaveSettings", args).Url
}

func (_ tHotels) ConfirmBooking(
		id int,
		booking interface{},
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	revel.Unbind(args, "booking", booking)
	return revel.MainRouter.Reverse("Hotels.ConfirmBooking", args).Url
}

func (_ tHotels) CancelBooking(
		id int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("Hotels.CancelBooking", args).Url
}

func (_ tHotels) Book(
		id int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "id", id)
	return revel.MainRouter.Reverse("Hotels.Book", args).Url
}


