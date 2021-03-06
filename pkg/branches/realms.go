package branches

import (
	"context"
	"sync"

	"github.com/brendoncarroll/go-state/cadata"
	"github.com/brendoncarroll/go-state/cells"
	"github.com/pkg/errors"
)

var (
	ErrNotExist = errors.New("volume does not exist")
	ErrExists   = errors.New("a volume already exists by that name")
)

func IsNotExist(err error) bool {
	return err == ErrNotExist
}

func IsExists(err error) bool {
	return err == ErrExists
}

// A Realm is a set of named branches.
type Realm interface {
	Get(ctx context.Context, name string) (*Branch, error)
	Create(ctx context.Context, name string) error
	Delete(ctx context.Context, name string) error
	ForEach(ctx context.Context, fn func(string) error) error
}

func CreateIfNotExists(ctx context.Context, r Realm, k string) error {
	if _, err := r.Get(ctx, k); err != nil {
		if IsNotExist(err) {
			return r.Create(ctx, k)
		}
		return err
	}
	return nil
}

type MemRealm struct {
	newStore func() cadata.Store
	newCell  func() cells.Cell

	mu       sync.RWMutex
	branches map[string]Branch
}

func NewMem(newStore func() cadata.Store, newCell func() cells.Cell) Realm {
	return &MemRealm{
		newStore: newStore,
		newCell:  newCell,
		branches: map[string]Branch{},
	}
}

func (r *MemRealm) Get(ctx context.Context, name string) (*Branch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	branch, exists := r.branches[name]
	branch.Annotations = copyAnotations(branch.Annotations)
	if !exists {
		return nil, ErrNotExist
	}
	return &branch, nil
}

func (r *MemRealm) Create(ctx context.Context, name string) error {
	if err := CheckName(name); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.branches[name]; exists {
		return ErrExists
	}
	r.branches[name] = Branch{
		Volume: Volume{
			Cell:     r.newCell(),
			VCStore:  r.newStore(),
			FSStore:  r.newStore(),
			RawStore: r.newStore(),
		},
	}
	return nil
}

func (r *MemRealm) Delete(ctx context.Context, name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.branches, name)
	return nil
}

func (r *MemRealm) ForEach(ctx context.Context, fn func(string) error) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for name := range r.branches {
		if err := fn(name); err != nil {
			return err
		}
	}
	return nil
}

func copyAnotations(x map[string]string) map[string]string {
	y := make(map[string]string, len(x))
	for k, v := range x {
		y[k] = v
	}
	return y
}
