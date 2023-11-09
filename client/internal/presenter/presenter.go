package presenter

import "github.com/jinzhu/copier"

type Presenter[M any, P any] struct {
	model *M
}

func NewPresenter[M any, P any](model *M) Presenter[M, P] {
	return Presenter[M, P]{model}
}

func (p Presenter[M, P]) Presenter() (*P, error) {
	var presenter P
	if err := copier.Copy(&presenter, p.model); err != nil {
		return nil, err
	}

	return &presenter, nil
}

func (p *Presenter[M, P]) PresenterFrom(model *M) (*P, error) {
	var presenter P
	if err := copier.Copy(&presenter, model); err != nil {
		return nil, err
	}

	return &presenter, nil
}
