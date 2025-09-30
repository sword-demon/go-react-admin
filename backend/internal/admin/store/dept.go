package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/sword-demon/go-react-admin/internal/pkg/model"
	"gorm.io/gorm"
)

// deptStore implements IDeptStore interface
type deptStore struct {
	db *gorm.DB
}

// newDeptStore creates a new dept store
func newDeptStore(db *gorm.DB) IDeptStore {
	return &deptStore{db: db}
}

// Create creates a new department
func (s *deptStore) Create(ctx context.Context, dept *model.Dept) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Build ancestors chain
		if dept.ParentID > 0 {
			var parent model.Dept
			if err := tx.First(&parent, dept.ParentID).Error; err != nil {
				return fmt.Errorf("parent dept not found: %w", err)
			}
			if parent.Ancestors == "" {
				dept.Ancestors = fmt.Sprintf("0,%d", parent.ID)
			} else {
				dept.Ancestors = fmt.Sprintf("%s,%d", parent.Ancestors, parent.ID)
			}
		} else {
			dept.Ancestors = "0"
		}

		return tx.Create(dept).Error
	})
}

// Update updates department information
func (s *deptStore) Update(ctx context.Context, dept *model.Dept) error {
	return s.db.WithContext(ctx).Model(dept).Updates(dept).Error
}

// Delete soft deletes a department
func (s *deptStore) Delete(ctx context.Context, id uint64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check if dept has children
		var count int64
		if err := tx.Model(&model.Dept{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("cannot delete department with children")
		}

		// Check if dept has users
		if err := tx.Model(&model.User{}).Where("dept_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("cannot delete department with users")
		}

		return tx.Delete(&model.Dept{}, id).Error
	})
}

// Get retrieves a department by ID
func (s *deptStore) Get(ctx context.Context, id uint64) (*model.Dept, error) {
	var dept model.Dept
	err := s.db.WithContext(ctx).
		Preload("Parent").
		First(&dept, id).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

// List retrieves all departments (for building tree)
func (s *deptStore) List(ctx context.Context) ([]*model.Dept, error) {
	var depts []*model.Dept
	err := s.db.WithContext(ctx).
		Where("status = ?", model.StatusEnabled).
		Order("order_num ASC").
		Find(&depts).Error
	if err != nil {
		return nil, err
	}
	return depts, nil
}

// GetChildren retrieves child departments
func (s *deptStore) GetChildren(ctx context.Context, parentID uint64) ([]*model.Dept, error) {
	var depts []*model.Dept

	// Build query based on parent ID
	query := s.db.WithContext(ctx).Model(&model.Dept{})

	if parentID == 0 {
		// Get all root departments
		query = query.Where("parent_id = ?", 0)
	} else {
		// Get all descendants using ancestors field
		var parent model.Dept
		if err := s.db.WithContext(ctx).First(&parent, parentID).Error; err != nil {
			return nil, err
		}

		// Find all children where ancestors contains parent ID
		ancestorsPattern := fmt.Sprintf("%%,%d,%%", parentID)
		query = query.Where(
			"parent_id = ? OR ancestors LIKE ? OR ancestors LIKE ? OR ancestors LIKE ?",
			parentID,
			ancestorsPattern,
			fmt.Sprintf("%d,%%", parentID),
			fmt.Sprintf("%%,%d", parentID),
		)
	}

	err := query.
		Where("status = ?", model.StatusEnabled).
		Order("order_num ASC").
		Find(&depts).Error

	if err != nil {
		return nil, err
	}
	return depts, nil
}

// GetDeptIDs retrieves all dept IDs including children (for data scope filtering)
func (s *deptStore) GetDeptIDs(ctx context.Context, deptID uint64, includeChildren bool) ([]uint64, error) {
	if !includeChildren {
		return []uint64{deptID}, nil
	}

	// Get department
	var dept model.Dept
	if err := s.db.WithContext(ctx).First(&dept, deptID).Error; err != nil {
		return nil, err
	}

	// Find all children departments
	var children []*model.Dept
	ancestorsPattern := fmt.Sprintf("%%,%d%%", deptID)
	err := s.db.WithContext(ctx).
		Model(&model.Dept{}).
		Where("id = ? OR ancestors LIKE ?", deptID, ancestorsPattern).
		Select("id").
		Find(&children).Error

	if err != nil {
		return nil, err
	}

	// Extract IDs
	ids := make([]uint64, 0, len(children))
	for _, child := range children {
		ids = append(ids, child.ID)
	}

	return ids, nil
}

// BuildTree builds department tree structure
func (s *deptStore) BuildTree(depts []*model.Dept) []*model.Dept {
	// Create map for quick lookup
	deptMap := make(map[uint64]*model.Dept)
	for _, dept := range depts {
		deptMap[dept.ID] = dept
	}

	// Build tree
	var roots []*model.Dept
	for _, dept := range depts {
		if dept.ParentID == 0 {
			roots = append(roots, dept)
		} else if parent, exists := deptMap[dept.ParentID]; exists {
			if parent.Children == nil {
				parent.Children = []*model.Dept{}
			}
			parent.Children = append(parent.Children, dept)
		}
	}

	return roots
}

// GetAncestorIDs parses ancestors string to ID slice
func GetAncestorIDs(ancestors string) []uint64 {
	if ancestors == "" || ancestors == "0" {
		return []uint64{}
	}

	parts := strings.Split(ancestors, ",")
	ids := make([]uint64, 0, len(parts))

	for _, part := range parts {
		var id uint64
		if _, err := fmt.Sscanf(part, "%d", &id); err == nil && id > 0 {
			ids = append(ids, id)
		}
	}

	return ids
}
