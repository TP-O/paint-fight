package presenter

import "github.com/jinzhu/copier"

type Presentation[M any, P any] struct {
	model *M
}

func NewPresenter[M any, P any](model *M) Presentation[M, P] {
	return Presentation[M, P]{model}
}

func (p Presentation[M, P]) Presenter() (*P, error) {
	return PresenterFrom[M, P](p.model)
}

func (p *Presentation[M, P]) PresenterFrom(model *M) (*P, error) {
	return PresenterFrom[M, P](model)
}

func PresenterFrom[F any, T any](from *F) (*T, error) {
	var to T
	if err := copier.Copy(&to, from); err != nil {
		return nil, err
	}

	return &to, nil
}
