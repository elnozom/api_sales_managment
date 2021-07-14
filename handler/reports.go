package handler

import (
	"fmt"
	"hand_held/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetStores(c echo.Context) error {
	var stores []model.Store
	// return c.JSON(http.StatusOK, "test")
	rows, err := h.db.Raw("EXEC GetStoreName").Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	defer rows.Close()
	for rows.Next() {
		var store model.Store
		rows.Scan(&store.StoreCode, &store.StoreName)
		stores = append(stores, store)
	}

	return c.JSON(http.StatusOK, stores)
}

func (h *Handler) GetAccount(c echo.Context) error {

	req := new(model.GetAccountRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	fmt.Println(req)

	var accounts []model.Account
	rows, err := h.db.Raw("EXEC GetAccount @Code = ?, @Name = ? , @Type = ?", req.Code, req.Name, req.Type).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var account model.Account
		rows.Scan(&account.Serial, &account.AccountCode, &account.AccountName)
		accounts = append(accounts, account)
	}

	return c.JSON(http.StatusOK, accounts)
}

func (h *Handler) GetItem(c echo.Context) error {

	req := new(model.GetItemRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	fmt.Println(req.BCode)

	var items []model.SingleItem
	rows, err := h.db.Raw("EXEC GetItemData @BCode = ?, @StoreCode = ?", req.BCode, req.StoreCode).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var item model.SingleItem
		err = rows.Scan(&item.Serial, &item.ItemName, &item.MinorPerMajor, &item.POSPP, &item.POSTP, &item.ByWeight)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		items = append(items, item)
	}

	return c.JSON(http.StatusOK, items)
}

func (h *Handler) GetItems(c echo.Context) error {

	type Req struct {
		StoreCode int `json:"StoreCode" validate:"required"`
	}

	req := new(Req)
	if err := c.Bind(req); err != nil {
		return err
	}
	fmt.Println(req.StoreCode)

	var items []model.Item
	rows, err := h.db.Raw("EXEC GetItems @StoreCode = ?", req.StoreCode).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var item model.Item
		err = rows.Scan(&item.Name, &item.Code, &item.Price, &item.Qnt, &item.AnQnt, &item.LimitedQnt, &item.StopSale, &item.PMin, &item.PMax)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		items = append(items, item)
	}

	return c.JSON(http.StatusOK, items)
}
