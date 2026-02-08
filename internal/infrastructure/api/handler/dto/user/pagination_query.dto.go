package dto

type Pagination struct {
	TempPage int `form:"page" binding:"omitempty,min=1"`
	TempSize int `form:"size" binding:"omitempty,min=1,max=100"`
}

const (
	defaultPage = 1
	defaultSize = 10
)

func (p *Pagination) Page() int {
	if p.TempPage == 0 {
		return defaultPage
	}
	return p.TempPage
}

func (p *Pagination) Size() int {
	if p.TempSize == 0 {
		return defaultSize
	}
	return p.TempSize
}
