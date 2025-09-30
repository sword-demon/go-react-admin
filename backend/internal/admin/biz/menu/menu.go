package menu

import (
	"context"

	"github.com/sword-demon/go-react-admin/internal/admin/biz"
	"github.com/sword-demon/go-react-admin/internal/admin/store"
	"github.com/sword-demon/go-react-admin/internal/pkg/model"
)

type menuBiz struct {
	store store.IStore
}

func NewMenuBiz(store store.IStore) biz.IMenuBiz {
	return &menuBiz{store: store}
}

func (b *menuBiz) Create(ctx context.Context, req *biz.CreateMenuRequest) (*biz.MenuResponse, error) {
	menu := &model.Menu{
		ParentID:  req.ParentID,
		MenuName:  req.MenuName,
		MenuType:  getMenuType(req.MenuType),
		Path:      req.Path,
		Component: req.Component,
		PermKey:   req.Perms,
		Icon:      req.Icon,
		OrderNum:  req.Sort,
		Status:    model.StatusEnabled,
		Visible:   model.StatusEnabled,
	}
	if err := b.store.Menus().Create(ctx, menu); err != nil {
		return nil, err
	}
	return b.toMenuResponse(menu), nil
}

func (b *menuBiz) Update(ctx context.Context, id uint64, req *biz.UpdateMenuRequest) error {
	menu, err := b.store.Menus().Get(ctx, id)
	if err != nil {
		return err
	}
	if req.MenuName != "" {
		menu.MenuName = req.MenuName
	}
	if req.Path != "" {
		menu.Path = req.Path
	}
	if req.Status >= 0 {
		menu.Status = uint8(req.Status)
	}
	return b.store.Menus().Update(ctx, menu)
}

func (b *menuBiz) Delete(ctx context.Context, id uint64) error {
	return b.store.Menus().Delete(ctx, id)
}

func (b *menuBiz) Get(ctx context.Context, id uint64) (*biz.MenuResponse, error) {
	menu, err := b.store.Menus().Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return b.toMenuResponse(menu), nil
}

func (b *menuBiz) GetTree(ctx context.Context) ([]*biz.MenuTreeNode, error) {
	menus, err := b.store.Menus().List(ctx)
	if err != nil {
		return nil, err
	}
	tree := b.store.Menus().BuildTree(menus)
	return b.toMenuTree(tree), nil
}

func (b *menuBiz) GetUserMenus(ctx context.Context, userID uint64) ([]*biz.MenuTreeNode, error) {
	menus, err := b.store.Menus().GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	tree := b.store.Menus().BuildTree(menus)
	return b.toMenuTree(tree), nil
}

func (b *menuBiz) toMenuResponse(menu *model.Menu) *biz.MenuResponse {
	return &biz.MenuResponse{
		ID:        menu.ID,
		MenuName:  menu.MenuName,
		ParentID:  menu.ParentID,
		MenuType:  string(menu.MenuType),
		Path:      menu.Path,
		Component: menu.Component,
		Perms:     menu.PermKey,
		Icon:      menu.Icon,
		Sort:      menu.OrderNum,
		Visible:   int8(menu.Visible),
		Status:    int8(menu.Status),
	}
}

func (b *menuBiz) toMenuTree(menus []*model.Menu) []*biz.MenuTreeNode {
	nodes := make([]*biz.MenuTreeNode, 0, len(menus))
	for _, menu := range menus {
		node := &biz.MenuTreeNode{
			MenuResponse: b.toMenuResponse(menu),
		}
		if len(menu.Children) > 0 {
			node.Children = b.toMenuTree(menu.Children)
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func getMenuType(menuType string) uint8 {
	switch menuType {
	case "D":
		return model.MenuTypeDirectory
	case "M":
		return model.MenuTypeMenu
	case "B":
		return model.MenuTypeButton
	default:
		return model.MenuTypeMenu
	}
}
