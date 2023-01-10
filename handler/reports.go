package handler

import (
	"fmt"
	"net/http"
	"sms/model"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *Handler) StockRpt(c echo.Context) error {

	var resp []model.IetmBalance
	var store int64
	if c.FormValue("store") != "" {
		store, _ = strconv.ParseInt(c.FormValue("store"), 0, 64)
	}
	rows, err := h.db.Raw("EXEC GetItemBalance @StoreCode = ? ", store).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var rec model.IetmBalance
		rows.Scan(&rec.Raseed, &rec.ItemName, &rec.RaseedReserved, &rec.AnRaseed, &rec.AnRaseedReserved, &rec.StoreCode, &rec.StoreName, &rec.DsiplayName)
		rec.RaseedNet = rec.Raseed - rec.RaseedReserved
		rec.AnRaseedNet = rec.AnRaseed - rec.AnRaseedReserved
		resp = append(resp, rec)
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetItemBalance(c echo.Context) error {

	var resp []model.IetmBalance
	rows, err := h.db.Raw("EXEC GetItemBalance @ItemSerial = ?", c.Param("serial")).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var rec model.IetmBalance
		rows.Scan(&rec.Raseed, &rec.ItemName, &rec.RaseedReserved, &rec.AnRaseed, &rec.AnRaseedReserved, &rec.StoreCode, &rec.StoreName, &rec.DsiplayName)
		rec.RaseedNet = rec.Raseed - rec.RaseedReserved
		rec.AnRaseedNet = rec.AnRaseed - rec.AnRaseedReserved
		resp = append(resp, rec)
	}

	return c.JSON(http.StatusOK, resp)
}
func (h *Handler) UpdateReservedForEmp(c echo.Context) error {
	code := c.Get("empCode")
	rows, err := h.db.Raw("EXEC UnReserveEmpDocs  @EmpCode = ? ", code).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	return c.JSON(http.StatusOK, "exited")
}
func (h *Handler) ValidateUser(c echo.Context) error {
	return c.JSON(http.StatusOK, "valid")
}

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

func (h *Handler) InsertOrder(c echo.Context) error {
	code := c.Get("empCode")
	req := new(model.InsertOrder)
	if err := c.Bind(req); err != nil {
		return err
	}
	orderNoRows, err := h.db.Raw("EXEC GetSalesOrderDocNo @StoreCode = ?, @TrSerial = ?", req.StoreCode, 30).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var orderNo int32
	for orderNoRows.Next() {
		orderNoRows.Scan(&orderNo)
	}
	orderNo = orderNo + 1
	fmt.Println(orderNo)

	rows, err := h.db.Raw("EXEC InsertTr05 @DocNo = ?, @StoreCode = ? , @EmpCode = ? , @AccountSerial =? ", orderNo, req.StoreCode, code, req.AccountSerial).Rows()
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
	var resp = new(model.InsertOrderResp)
	resp.Serial = serial
	resp.No = int(orderNo)
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) InsertOrderItem(c echo.Context) error {
	req := new(model.InsertOrderItem)
	if err := c.Bind(req); err != nil {
		return err
	}
	rows, err := h.db.Raw("EXEC InsertTr06 @HeadSerial = ?, @ItemSerial = ? , @Qnt = ? , @Price = ? , @QntAntherUnit = ? , @PriceMax = ? , @PriceMin = ? , @MinorPerMajor = ? , @Branch = ?", req.HeadSerial, req.ItemSerial, req.Qnt, req.Price, req.QntAntherUnit, req.PriceMax, req.PriceMin, req.MinorPerMajor, req.Branch).Rows()
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
		err = rows.Scan(&item.Serial, &item.ItemHaveAntherUnit, &item.BarCode, &item.ItemName, &item.ItemSerial, &item.Qnt, &item.QntAntherUnit, &item.AvgWeight, &item.Price, &item.PriceMax, &item.PriceMin, &item.Total, &item.StoreCode, &item.StoreName)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		items = append(items, item)
	}

	totalsRows, err := h.db.Raw("EXEC StkTr06CalculateTotals @HeadSerial = ?", req.Serial).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer totalsRows.Close()
	var totals model.OrderTotals
	for totalsRows.Next() {
		err = totalsRows.Scan(&totals.TotalPackages, &totals.TotalCash)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	response := model.GetOrderItemsResonse{
		Items:  items,
		Totals: totals,
	}

	return c.JSON(http.StatusOK, response)
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
		err = rows.Scan(&item.Serial, &item.Name, &item.Code, &item.Price, &item.Qnt, &item.AnQnt, &item.ItemHaveAntherUnit, &item.LimitedQnt, &item.StopSale, &item.PMax, &item.PMin, &item.AvrWeight, &item.MinorPerMajor)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		items = append(items, item)
	}

	return c.JSON(http.StatusOK, items)
}

func (h *Handler) UpdateItem(c echo.Context) error {
	type Req struct {
		LQvalue  *float64
		STValue  *bool
		MinValue *float64
		MaxValue *float64
		Serial   *int
		Branch   *int
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

func (h *Handler) UpdateOrder(c echo.Context) error {
	type Req struct {
		Reserved    *bool
		AuditCode   *int
		DeliveryFee float64
		DriverName  *string
	}

	req := new(Req)
	if err := c.Bind(req); err != nil {
		return err
	}
	rows, err := h.db.Raw("EXEC StkTr05Update  @Serial = ? ,@Reserved = ?,@AuditCode =? ,@DeliveryFee =? , @DriverName = ?", c.Param("Serial"), req.Reserved, req.AuditCode, req.DeliveryFee, req.DriverName).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	return c.JSON(http.StatusOK, "exited")
}

func (h *Handler) ListOrders(c echo.Context) error {
	code := c.Get("empCode")
	type Req struct {
		Finished bool
	}

	req := new(Req)
	if err := c.Bind(req); err != nil {
		return err
	}
	fmt.Println("req")
	fmt.Println(req)
	rows, err := h.db.Raw("EXEC ListTr05 @EmpCode = ? ,   @Finished = ? ", code, req.Finished).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var orders []model.Order
	defer rows.Close()
	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.Serial, &order.StkTr01Serial, &order.DocNo, &order.DocDate, &order.EmpCode, &order.TotalCash, &order.EmpName, &order.CustomerName, &order.CustomerCode, &order.CustomerSerial, &order.Reserved, &order.DriverName, &order.DeliveryFee, &order.Finished)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		order.DocDate = strings.Replace(order.DocDate, "T", " ", 1)
		order.DocDate = strings.Replace(order.DocDate, "Z", " ", 1)
		orders = append(orders, order)
	}

	return c.JSON(http.StatusOK, orders)
}

func (h *Handler) ListStoreOrders(c echo.Context) error {
	code := c.Get("empCode")
	var employee model.Emp
	type Req struct {
		Finished bool
	}

	req := new(Req)

	if err := c.Bind(req); err != nil {
		return err
	}
	err := h.db.Raw("EXEC GetEmp @EmpCode = ?  ", code).Row().Scan(
		&employee.EmpName, &employee.EmpPassword, &employee.EmpCode, &employee.SecLevel, &employee.FixEmpStore)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	rows, err := h.db.Raw("EXEC STCOrderList @StoreCode = ? ,  @Finished = ?", employee.FixEmpStore, req.Finished).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var orders []model.Order
	defer rows.Close()
	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.Serial, &order.StcEmpName, &order.DocNo, &order.DocDate, &order.StoreName, &order.TotalCash, &order.DriverName, &order.DeliveryFee, &order.CustomerName, &order.CustomerCode, &order.EmpName, &order.SalesOrderNo, &order.Finished)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		order.DocDate = strings.Replace(order.DocDate, "T", " ", 1)
		order.DocDate = strings.Replace(order.DocDate, "Z", " ", 1)
		orders = append(orders, order)
	}

	return c.JSON(http.StatusOK, orders)
}

func (h *Handler) ListStoreOrderItems(c echo.Context) error {
	rows, err := h.db.Raw("EXEC STCDetailsList @Serial = ? ", c.Param("Serial")).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var orders []model.OrderItem
	defer rows.Close()
	for rows.Next() {
		var order model.OrderItem
		err = rows.Scan(&order.Serial, &order.BarCode, &order.ItemName, &order.ItemSerial, &order.Qnt, &order.QntAntherUnit, &order.AvgWeight, &order.Price, &order.Total)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		orders = append(orders, order)
	}

	return c.JSON(http.StatusOK, orders)
}

func (h *Handler) CloseStoreOrderItems(c echo.Context) error {

	var resp int
	err := h.db.Raw("EXEC STCOrdersClose @Serial = ? , @StcEmp = ? ", c.Param("serial"), c.Param("emp")).Row().Scan(&resp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}
func (h *Handler) UpdateStoreOrderItems(c echo.Context) error {
	type Req struct {
		Serial int
		Qnt    float64
	}

	req := new(Req)
	var resp int
	if err := c.Bind(req); err != nil {
		return err
	}
	err := h.db.Raw("EXEC STCDetailsUpdate @Serial = ? , @Qnt = ? ", req.Serial, req.Qnt).Row().Scan(&resp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateOrderItem(c echo.Context) error {
	type Req struct {
		Serial        int
		Qnt           float64
		Price         float64
		QntAntherUnit float64
		Branch        int
	}

	req := new(Req)
	if err := c.Bind(req); err != nil {
		return err
	}
	rows, err := h.db.Raw("EXEC UpdateTr06  @Serial = ? , @Price = ? , @Qnt = ? , @Branch = ? , @QntAntherUnit = ? ", req.Serial, req.Price, req.Qnt, req.Branch, req.QntAntherUnit).Rows()
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
		Serial int
		DocNo  int
	}
	code := c.Get("empCode")
	req := new(Req)
	if err := c.Bind(req); err != nil {
		return err
	}

	_, err := h.db.Raw("EXEC StkTr05Update  @Serial = ? ,@UpdateTotalCash = ? , @AuditCode = ?", req.Serial, 1, code).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	_, err = h.db.Raw("EXEC StcOrdersInsert  @HeadSerial = ?", req.Serial).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	rows, err := h.db.Raw("EXEC STCOrderList @HeadNo = ? ", req.DocNo).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var orders []model.Order
	var details [][]model.OrderItem
	defer rows.Close()
	for rows.Next() {
		var d []model.OrderItem
		var order model.Order
		err = rows.Scan(&order.Serial, &order.StcEmpName, &order.DocNo, &order.DocDate, &order.StoreName, &order.TotalCash, &order.DriverName, &order.DeliveryFee, &order.CustomerName, &order.CustomerCode, &order.EmpName, &order.SalesOrderNo, &order.Finished)
		// err = rows.Scan(&order.Serial, &order.DocNo, &order.DocDate, &order.StoreName, &order.TotalCash, &order.DriverName, &order.DeliveryFee, &order.CustomerName, &order.CustomerCode, &order.EmpName, &order.SalesOrderNo, &order.Finished)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		order.DocDate = strings.Replace(order.DocDate, "T", " ", 1)
		order.DocDate = strings.Replace(order.DocDate, "Z", " ", 1)
		orders = append(orders, order)

		// get items fior current head

		detailsRows, err := h.db.Raw("EXEC STCDetailsList @Serial = ? ", order.Serial).Rows()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		defer rows.Close()
		for detailsRows.Next() {
			var item model.OrderItem
			err = detailsRows.Scan(&item.Serial, &item.BarCode, &item.ItemName, &item.ItemSerial, &item.Qnt, &item.QntAntherUnit, &item.AvgWeight, &item.Price, &item.Total)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			item.StoreName = order.StoreName
			d = append(d, item)
		}

		details = append(details, d)

	}
	return c.JSON(http.StatusOK, details)
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

func (h *Handler) Login(c echo.Context) error {
	req := new(model.LoginReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	code, err := strconv.ParseUint(req.EmpCode, 10, 32)

	if err != nil {
		return c.JSON(http.StatusBadRequest, "empcode_not_valid"+err.Error())
	}
	var employee model.Emp
	rows, err := h.db.Raw("EXEC GetEmp @EmpCode = ?", code).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var item model.Emp
		err = rows.Scan(&item.EmpName, &item.EmpPassword, &item.EmpCode, &item.SecLevel, &item.FixEmpStore)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		employee = item
	}

	if employee.EmpPassword == "" {
		return c.JSON(http.StatusBadGateway, "incorrect_empcode")

	}
	fmt.Println(employee.EmpPassword)
	if employee.EmpPassword != req.EmpPassword {
		return c.JSON(http.StatusBadGateway, "incorrect_password")

	}

	accessToken, err := h.tokenMaker.CreateToken(
		uint32(code),
		time.Duration(99*1000000000000),
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	response := model.LoginResponse{accessToken, employee}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetEmp(c echo.Context) error {
	code := c.Get("empCode")
	var employee model.Emp
	rows, err := h.db.Raw("EXEC GetEmp @EmpCode = ?", code).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var item model.Emp
		err = rows.Scan(&item.EmpName, &item.EmpPassword, &item.EmpCode, &item.SecLevel, &item.FixEmpStore)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		employee = item
	}

	return c.JSON(http.StatusOK, employee)
}
