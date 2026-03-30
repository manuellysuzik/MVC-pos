package pedidosvc

import (
	"errors"
	"fmt"

	"ecommerce-api/internal/domain/entity"
)

type PedidoService struct {
	repo PedidoRepository
}

func New(repo PedidoRepository) *PedidoService {
	return &PedidoService{repo: repo}
}

func (s *PedidoService) Create(p *entity.Pedido) error {
	if p.ClienteID == 0 {
		return entity.NewValidationError("cliente_id é obrigatório")
	}
	if len(p.Itens) == 0 {
		return entity.NewValidationError("pedido deve ter pelo menos um item")
	}
	total := 0.0
	for _, item := range p.Itens {
		if item.Quantidade <= 0 {
			return entity.NewValidationError("quantidade deve ser maior que zero")
		}
		if item.PrecoUnitario <= 0 {
			return entity.NewValidationError("preco_unitario deve ser maior que zero")
		}
		total += float64(item.Quantidade) * item.PrecoUnitario
	}
	p.Total = total
	if p.Status == "" {
		p.Status = "pendente"
	}
	if err := s.repo.Create(p); err != nil {
		return fmt.Errorf("criando pedido: %w", err)
	}
	return nil
}

func (s *PedidoService) FindAll() ([]entity.Pedido, error) {
	pedidos, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("buscando pedidos: %w", err)
	}
	return pedidos, nil
}

func (s *PedidoService) FindByID(id uint) (*entity.Pedido, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return nil, entity.ErrNotFound
		}
		return nil, fmt.Errorf("buscando pedido %d: %w", id, err)
	}
	return p, nil
}

func (s *PedidoService) FindByStatus(status string) ([]entity.Pedido, error) {
	if status == "" {
		return nil, entity.NewValidationError("status é obrigatório para busca")
	}
	pedidos, err := s.repo.FindByStatus(status)
	if err != nil {
		return nil, fmt.Errorf("buscando pedidos por status: %w", err)
	}
	return pedidos, nil
}

func (s *PedidoService) FindByClienteID(clienteID uint) ([]entity.Pedido, error) {
	pedidos, err := s.repo.FindByClienteID(clienteID)
	if err != nil {
		return nil, fmt.Errorf("buscando pedidos do cliente %d: %w", clienteID, err)
	}
	return pedidos, nil
}

func (s *PedidoService) Count() (int64, error) {
	count, err := s.repo.Count()
	if err != nil {
		return 0, fmt.Errorf("contando pedidos: %w", err)
	}
	return count, nil
}

func (s *PedidoService) Update(p *entity.Pedido) error {
	if p.ID == 0 {
		return entity.NewValidationError("id é obrigatório para atualização")
	}
	if _, err := s.repo.FindByID(p.ID); err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return entity.ErrNotFound
		}
		return fmt.Errorf("verificando pedido %d: %w", p.ID, err)
	}
	if err := s.repo.Update(p); err != nil {
		return fmt.Errorf("atualizando pedido %d: %w", p.ID, err)
	}
	return nil
}

func (s *PedidoService) Delete(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return entity.ErrNotFound
		}
		return fmt.Errorf("verificando pedido %d: %w", id, err)
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("deletando pedido %d: %w", id, err)
	}
	return nil
}
