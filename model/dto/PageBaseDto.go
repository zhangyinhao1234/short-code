package dto

type PageBaseDto struct {
	Total int64 `json:"total"`
}

type PageQueryParams struct {
	Size    int `json:"size" validate:"required,min=1" label:"每页条数" `
	PageNum int `json:"pageNum" validate:"required,min=1" label:"第几页"`
}

func (p *PageQueryParams) GetLimit() int {
	return p.Size
}

func (p *PageQueryParams) GetOffSet() int {
	return (p.PageNum - 1) * p.Size
}
