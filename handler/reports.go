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

func (h *Handler) GetEmp(c echo.Context) error {
	type Req struct {
		EmpCode int
	}
	req := new(model.EmpReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	fmt.Println(req.EmpCode)

	var employee []model.Emp
	rows, err := h.db.Raw("EXEC GetEmp @EmpCode = ?", req.EmpCode).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var item model.Emp
		err = rows.Scan(&item.EmpName, &item.EmpPassword, &item.EmpCode)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		employee = append(employee, item)
	}

	return c.JSON(http.StatusOK, employee[0])
}

func (h *Handler) InsertOrder(c echo.Context) error {

	req := new(model.InsertOrder)
	if err := c.Bind(req); err != nil {
		return err
	}
	rows, err := h.db.Raw("EXEC InsertTr05 @DocNo = ?, @StoreCode = ? , @EmpCode = ? , @AccountSerial =? ", req.DocNo, req.StoreCode, req.EmpCode, req.AccountSerial).Rows()
	defer rows.Close()
	var serial int
	for rows.Next() {
		err = rows.Scan(&serial)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, serial)
}

func (h *Handler) InsertOrderItem(c echo.Context) error {

	req := new(model.InsertOrderItem)
	if err := c.Bind(req); err != nil {
		return err
	}
	rows, err := h.db.Raw("EXEC InsertTr06 @HeadSerial = ?, @ItemSerial = ? , @Qnt = ? , @Price = ?", req.HeadSerial, req.ItemSerial, req.Qnt, req.Price).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	type Resp struct {
		TotalPackages int
		TotalCash     float64
		Serial        int32
	}
	var resp = new(Resp)
	for rows.Next() {
		err = rows.Scan(&resp.TotalPackages, &resp.TotalCash)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	if rows.NextResultSet() {
		for rows.Next() {
			fmt.Println("asdasdasd")
			err = rows.Scan(&resp.Serial)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
		}

	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetItem(c echo.Context) error {
	req := new(model.GetItemRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	fmt.Println(req.BCode)
	var items []model.SingleItem
	rows, err := h.db.Raw("EXEC GetItemData @BCode = ?, @StoreCode = ?, @Name = ? ", req.BCode, req.StoreCode, req.Name).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var item model.SingleItem
		err = rows.Scan(&item.Serial, &item.ItemName, &item.MinorPerMajor, &item.POSPP, &item.POSTP, &item.ByWeight, &item.WithExp, &item.ItemHasAntherUnit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		items = append(items, item)
	}

	return c.JSON(http.StatusOK, items)
}

func (h *Handler) GetSalesOrderDocNo(c echo.Context) error {
	type Req struct {
		StoreCode int `json:"StoreCode" validate:"required"`
	}

	req := new(Req)

	if err := c.Bind(req); err != nil {
		return err
	}
	fmt.Println("StoreCode")
	fmt.Println(req.StoreCode)
	rows, err := h.db.Raw("EXEC GetSalesOrderDocNo @StoreCode = ?, @TrSerial = ?", req.StoreCode, 30).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var resp int32
	for rows.Next() {
		rows.Scan(&resp)
	}
	return c.JSON(http.StatusOK, resp+1)
}

func (h *Handler) GetOrderItems(c echo.Context) error {

	req := new(model.GetOrderItemsRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	fmt.Println(req.Serial)

	var items []model.OrderItem
	rows, err := h.db.Raw("EXEC GetDocItemData @Serial = ?", req.Serial).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var item model.OrderItem
		err = rows.Scan(&item.Serial, &item.BarCode, &item.ItemName, &item.Qnt, &item.Price, &item.Total)
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
		err = rows.Scan(&item.Serial, &item.Name, &item.Code, &item.Price, &item.Qnt, &item.AnQnt, &item.LimitedQnt, &item.StopSale, &item.PMax, &item.PMin)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		items = append(items, item)
	}

	return c.JSON(http.StatusOK, items)
}

func (h *Handler) UpdateItem(c echo.Context) error {
	type Req struct {
		LQvalue  *bool
		STValue  *bool
		MinValue *float64
		MaxValue *float64
		Serial   *int
	}

	req := new(Req)
	if err := c.Bind(req); err != nil {
		return err
	}

	fmt.Println(req.MinValue)

	_, err := h.db.Raw("EXEC UpdateMs  @LQvalue = ? , @STValue = ? , @MinValue = ? , @MaxValue = ? , @Serial = ? ", req.LQvalue, req.STValue, req.MinValue, req.MaxValue, req.Serial).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "updated")
}

func (h *Handler) ListOrders(c echo.Context) error {
	rows, err := h.db.Raw("EXEC ListTr05").Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var orders []model.Order
	defer rows.Close()
	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.DocNo, &order.DocDate, &order.EmpCode, &order.TotalCash, &order.EmpName)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		orders = append(orders, order)
	}

	return c.JSON(http.StatusOK, orders)
}

func (h *Handler) UpdateOrderItem(c echo.Context) error {
	type Req struct {
		Serial int
		Qnt    float64
	}

	req := new(Req)
	if err := c.Bind(req); err != nil {
		return err
	}
	rows, err := h.db.Raw("EXEC UpdateTr06  @Serial = ? , @Qnt = ? ", req.Serial, req.Qnt).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	var item model.OrderTotals
	for rows.Next() {
		err = rows.Scan(&item.TotalPackages, &item.TotalCash)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, item)
}
func (h *Handler) CloseOrder(c echo.Context) error {
	type Req struct {
		Serial    int
		TotalCash float64
	}

	req := new(Req)
	if err := c.Bind(req); err != nil {
		return err
	}

	_, err := h.db.Raw("EXEC CloseTr05  @Serial = ? , @TotalCash = ? ", req.Serial, req.TotalCash).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "closed")
}

func (h *Handler) DeleteOrderItem(c echo.Context) error {
	type Req struct {
		Serial    int
		TotalCash float64
	}

	req := new(Req)
	if err := c.Bind(req); err != nil {
		return err
	}

	rows, err := h.db.Raw("EXEC DeleteTr06  @Serial = ?", req.Serial).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	var item model.OrderTotals
	for rows.Next() {
		err = rows.Scan(&item.TotalPackages, &item.TotalCash)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, item)
}
