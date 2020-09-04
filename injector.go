// +build wireinject

package acc

import (
	"github.com/google/wire"
)

func InitSubmitter() Submitter {
	wire.Build(
		NewSeleniumSubmitter,
		wire.Bind(new(Submitter), new(*SeleniumSubmitter)),
	)
	return nil
}
