package dept

import (
	"context"

	"github.com/sword-demon/go-react-admin/internal/admin/store"
	"github.com/sword-demon/go-react-admin/internal/pkg/model"
)

type deptBiz struct {
	store store.IStore
}

// Keep interface here to avoid import cycle
type IDeptBiz interface {
	Create(ctx context.Context, req *CreateDeptRequest) (*DeptResponse, error)
	Update(ctx context.Context, id uint64, req *UpdateDeptRequest) error
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, id uint64) (*DeptResponse, error)
	GetTree(ctx context.Context) ([]*DeptTreeNode, error)
}

type CreateDeptRequest struct {
	ParentID uint64 `json:"parent_id"`
	DeptName string `json:"dept_name" binding:"required"`
	Leader   string `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Sort     int    `json:"sort"`
}

type UpdateDeptRequest struct {
	DeptName string `json:"dept_name"`
	Leader   string `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Sort     int    `json:"sort"`
	Status   int8   `json:"status"`
}

type DeptResponse struct {
	ID       uint64 `json:"id"`
	ParentID uint64 `json:"parent_id"`
	DeptName string `json:"dept_name"`
	Leader   string `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Sort     int    `json:"sort"`
	Status   int8   `json:"status"`
}

type DeptTreeNode struct {
	*DeptResponse
	Children []*DeptTreeNode `json:"children,omitempty"`
}

func NewDeptBiz(store store.IStore) IDeptBiz {
	return &deptBiz{store: store}
}

func (b *deptBiz) Create(ctx context.Context, req *CreateDeptRequest) (*DeptResponse, error) {
	dept := &model.Dept{
		ParentID: req.ParentID,
		DeptName: req.DeptName,
		Leader:   req.Leader,
		Phone:    req.Phone,
		Email:    req.Email,
		OrderNum: req.Sort,
		Status:   model.StatusEnabled,
	}
	if err := b.store.Depts().Create(ctx, dept); err != nil {
		return nil, err
	}
	return b.toDeptResponse(dept), nil
}

func (b *deptBiz) Update(ctx context.Context, id uint64, req *UpdateDeptRequest) error {
	dept, err := b.store.Depts().Get(ctx, id)
	if err != nil {
		return err
	}
	if req.DeptName != "" {
		dept.DeptName = req.DeptName
	}
	if req.Leader != "" {
		dept.Leader = req.Leader
	}
	if req.Status >= 0 {
		dept.Status = uint8(req.Status)
	}
	return b.store.Depts().Update(ctx, dept)
}

func (b *deptBiz) Delete(ctx context.Context, id uint64) error {
	return b.store.Depts().Delete(ctx, id)
}

func (b *deptBiz) Get(ctx context.Context, id uint64) (*DeptResponse, error) {
	dept, err := b.store.Depts().Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return b.toDeptResponse(dept), nil
}

func (b *deptBiz) GetTree(ctx context.Context) ([]*DeptTreeNode, error) {
	depts, err := b.store.Depts().List(ctx)
	if err != nil {
		return nil, err
	}
	tree := b.store.Depts().BuildTree(depts)
	return b.toDeptTree(tree), nil
}

func (b *deptBiz) toDeptResponse(dept *model.Dept) *DeptResponse {
	return &DeptResponse{
		ID:       dept.ID,
		ParentID: dept.ParentID,
		DeptName: dept.DeptName,
		Leader:   dept.Leader,
		Phone:    dept.Phone,
		Email:    dept.Email,
		Sort:     dept.OrderNum,
		Status:   int8(dept.Status),
	}
}

func (b *deptBiz) toDeptTree(depts []*model.Dept) []*DeptTreeNode {
	nodes := make([]*DeptTreeNode, 0, len(depts))
	for _, dept := range depts {
		node := &DeptTreeNode{
			DeptResponse: b.toDeptResponse(dept),
		}
		if len(dept.Children) > 0 {
			node.Children = b.toDeptTree(dept.Children)
		}
		nodes = append(nodes, node)
	}
	return nodes
}
