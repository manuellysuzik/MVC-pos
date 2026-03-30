package clientesvc

import (
	"errors"
	"fmt"

	"ecommerce-api/internal/domain/entity"
)

type ClienteService struct {
	repo ClienteRepository
}

func New(repo ClienteRepository) *ClienteService {
	return &ClienteService{repo: repo}
}

func (s *ClienteService) Create(c *entity.Cliente) error {
	if c.Nome == "" {
		return entity.NewValidationError("nome é obrigatório")
	}
	if c.Email == "" {
		return entity.NewValidationError("email é obrigatório")
	}
	if err := s.repo.Create(c); err != nil {
		return fmt.Errorf("criando cliente: %w", err)
	}
	return nil
}

func (s *ClienteService) FindAll() ([]entity.Cliente, error) {
	clientes, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("buscando clientes: %w", err)
	}
	return clientes, nil
}

func (s *ClienteService) FindByID(id uint) (*entity.Cliente, error) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return nil, entity.ErrNotFound
		}
		return nil, fmt.Errorf("buscando cliente %d: %w", id, err)
	}
	return c, nil
}

func (s *ClienteService) FindByName(nome string) ([]entity.Cliente, error) {
	if nome == "" {
		return nil, entity.NewValidationError("nome é obrigatório para busca")
	}
	clientes, err := s.repo.FindByName(nome)
	if err != nil {
		return nil, fmt.Errorf("buscando clientes por nome: %w", err)
	}
	return clientes, nil
}

func (s *ClienteService) Count() (int64, error) {
	count, err := s.repo.Count()
	if err != nil {
		return 0, fmt.Errorf("contando clientes: %w", err)
	}
	return count, nil
}

func (s *ClienteService) Update(c *entity.Cliente) error {
	if c.ID == 0 {
		return entity.NewValidationError("id é obrigatório para atualização")
	}
	if _, err := s.repo.FindByID(c.ID); err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return entity.ErrNotFound
		}
		return fmt.Errorf("verificando cliente %d: %w", c.ID, err)
	}
	if err := s.repo.Update(c); err != nil {
		return fmt.Errorf("atualizando cliente %d: %w", c.ID, err)
	}
	return nil
}

func (s *ClienteService) Delete(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return entity.ErrNotFound
		}
		return fmt.Errorf("verificando cliente %d: %w", id, err)
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("deletando cliente %d: %w", id, err)
	}
	return nil
}
