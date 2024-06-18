package choreography

import (
	"short-code/model/do"
	"short-code/service/domain/binding"
	"short-code/service/domain/unuse"
)

var (
	bindingService   = binding.BindingService{}
	unUseCodeService = unuse.UnUseCodeService{}
)

type BindDataService struct {
}

func (e *BindDataService) Bind(data *string) (*string, *do.ShortCodeError) {
	shotCode, shotErr := unUseCodeService.Poll()
	if shotErr != nil {
		return nil, shotErr
	}
	shotCode, shotErr = bindingService.Binding(shotCode, data)
	if shotErr != nil {
		return nil, shotErr
	}
	return shotCode, nil
}

func (e *BindDataService) Flush() {
	bindingService.Flush()
}

func (e *BindDataService) LoadBindCacheInLocal() {
	go func() {
		bindingService.LoadBindCacheInLocal()
	}()
}
