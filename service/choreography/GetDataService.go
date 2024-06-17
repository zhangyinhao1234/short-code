package choreography

import "short-code/model/do"

type GetDataService struct {
}

func (e *GetDataService) GetByShotCode(shotCode *string) (*string, *do.ShortCodeError) {
	return bindingService.GetByShotCode(shotCode)
}
