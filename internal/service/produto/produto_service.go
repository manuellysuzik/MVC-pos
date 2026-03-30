package produtosvc

import (
	"errors"
	"fmt"

	"ecommerce-api/internal/domain/entity"
)

type ProdutoService struct {
	repo ProdutoRepository
}

func New(repo ProdutoRepository) *ProdutoService {
	return &ProdutoService{repo: repo}
}

func (s *ProdutoService) Create(p *entity.Produto) error {
	if p.Nome == "" {
		return entity.NewValidationError("nome é obrigatório")
	}
	if p.Preco <= 0 {
		return entity.NewValidationError("preço deve ser maior que zero")
	}
	if err := s.repo.Create(p); err != nil {
		return fmt.Errorf("criando produto: %w", err)
	}
	return nil
}

func (s *ProdutoService) FindAll() ([]entity.Produto, error) {
	produtos, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("buscando produtos: %w", err)
	}
	return produtos, nil
}

func (s *ProdutoService) FindByID(id uint) (*entity.Produto, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return nil, entity.ErrNotFound
		}
		return nil, fmt.Errorf("buscando produto %d: %w", id, err)
	}
	return p, nil
}

func (s *ProdutoService) FindByName(nome string) ([]entity.Produto, error) {
	if nome == "" {
		return nil, entity.NewValidationError("nome é obrigatório para busca")
	}
	produtos, err := s.repo.FindByName(nome)
	if err != nil {
		return nil, fmt.Errorf("buscando produtos por nome: %w", err)
	}
	return produtos, nil
}

func (s *ProdutoService) Count() (int64, error) {
	count, err := s.repo.Count()
	if err != nil {
		return 0, fmt.Errorf("contando produtos: %w", err)
	}
	return count, nil
}

func (s *ProdutoService) Update(p *entity.Produto) error {
	if p.ID == 0 {
		return entity.NewValidationError("id é obrigatório para atualização")
	}
	if _, err := s.repo.FindByID(p.ID); err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return entity.ErrNotFound
		}
		return fmt.Errorf("verificando produto %d: %w", p.ID, err)
	}
	if err := s.repo.Update(p); err != nil {
		return fmt.Errorf("atualizando produto %d: %w", p.ID, err)
	}
	return nil
}

func (s *ProdutoService) Delete(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return entity.ErrNotFound
		}
		return fmt.Errorf("verificando produto %d: %w", id, err)
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("deletando produto %d: %w", id, err)
	}
	return nil
}
