// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package acc

// Injectors from injector.go:

func InitSubmitter() Submitter {
	seleniumSubmitter := NewSeleniumSubmitter()
	return seleniumSubmitter
}