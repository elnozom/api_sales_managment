package handler

import (
	"net/http"
	"sms/model"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *Handler) PosOptionsRead(c echo.Context) error {
	req := new(model.InsertOrder)
	if err := c.Bind(req); err != nil {
		return err
	}
	var resp model.PostOptions
	err := h.db.Raw("EXEC PosOptionsRead").Row().Scan(
		&resp.ReportTitle,
		&resp.ReportPhone,
		&resp.BonMsg3,
		&resp.BonMsg4,
		&resp.BonMsg5,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}
func (h *Handler) InvoiceOrderDocNo(c echo.Context) error {
	req := new(model.InsertOrder)
	if err := c.Bind(req); err != nil {
		return err
	}
	orderNoRows, err := h.db.Raw("EXEC StkTrInvoiceDocNo @TrSerial = ?", 30).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var orderNo int32
	for orderNoRows.Next() {
		orderNoRows.Scan(&orderNo)
	}
	orderNo = orderNo + 1
	return c.JSON(http.StatusOK, orderNo)

}

func (h *Handler) InsertInvoiceOrder(c echo.Context) error {
	code := c.Get("empCode")
	req := new(model.InsertOrder)
	if err := c.Bind(req); err != nil {
		return err
	}
	orderNoRows, err := h.db.Raw("EXEC StkTrInvoiceDocNo @TrSerial = ?", 30).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var orderNo int32
	for orderNoRows.Next() {
		orderNoRows.Scan(&orderNo)
	}

	orderNo = orderNo + 1
	rows, err := h.db.Raw("EXEC StkTrInvoiceHeadInsert @DocNo = ?, @StoreCode = ? , @EmpCode = ? , @AccountSerial =? ", orderNo, req.StoreCode, code, req.AccountSerial).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
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

func (h *Handler) InsertInvoiceOrderItem(c echo.Context) error {
	req := new(model.InsertOrderItem)
	if err := c.Bind(req); err != nil {
		return err
	}
	var serial int
	err := h.db.Raw("EXEC StkTrInvoiceDetailsInsert @HeadSerial = ?, @ItemSerial = ?, @Qnt = ? , @Price = ? , @QntAntherUnit = ? , @PriceMax = ? , @PriceMin = ? , @MinorPerMajor = ? , @Branch = ?", req.HeadSerial, req.ItemSerial, req.Qnt, req.Price, req.QntAntherUnit, req.PriceMax, req.PriceMin, req.MinorPerMajor, req.Branch).Row().Scan(&serial)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, serial)
}

func (h *Handler) DeleteInvoiceOrder(c echo.Context) error {
	rows, err := h.db.Raw("EXEC StkTrInvoiceHeadDelete  @Serial = ?", c.Param("serial")).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	return c.JSON(http.StatusOK, "deleted")
}
func (h *Handler) ListInvoiceOrders(c echo.Context) error {
	code := c.Get("empCode")
	type Req struct {
		Finished bool
		Deleted  *bool
		DateFrom string
		DateTo   string
	}

	req := new(Req)
	if err := c.Bind(req); err != nil {
		return err
	}
	rows, err := h.db.Raw("EXEC StkTrInvoiceHeadList @EmpCode = ? ,   @Finished = ? , @Deleted = ? , @DateFrom = ? , @DateTo = ? ", code, req.Finished, req.Deleted, req.DateFrom, req.DateTo).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var orders []model.Order
	defer rows.Close()
	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.Serial, &order.DocNo, &order.DocDate, &order.EmpCode, &order.TotalCash, &order.EmpName, &order.CustomerName, &order.CustomerCode, &order.CustomerSerial, &order.Reserved, &order.Finished, &order.Deleted)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		order.DocDate = strings.Replace(order.DocDate, "T", " ", 1)
		order.DocDate = strings.Replace(order.DocDate, "Z", " ", 1)
		orders = append(orders, order)
	}

	return c.JSON(http.StatusOK, orders)
}

func (h *Handler) ListInvoiceItems(c echo.Context) error {
	var items []model.OrderItem
	rows, err := h.db.Raw("EXEC StkTrInvoiceDetailsList @Serial = ?", c.Param("serial")).Rows()
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

	return c.JSON(http.StatusOK, items)
}
func (h *Handler) DeleteInvoiceOrderItem(c echo.Context) error {
	rows, err := h.db.Raw("EXEC StkTrInvoiceDetailsDelete  @Serial = ?", c.Param("serial")).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	return c.JSON(http.StatusOK, "deleted")
}
func (h *Handler) UpdateInvoiceOrder(c echo.Context) error {
	type Req struct {
		Reserved  *bool
		AuditCode *int
	}

	req := new(Req)
	if err := c.Bind(req); err != nil {
		return err
	}
	rows, err := h.db.Raw("EXEC StkTrInvoiceHeadUpdate  @Serial = ? ,@Reserved = ?,@AuditCode =?", c.Param("Serial"), req.Reserved, req.AuditCode).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	return c.JSON(http.StatusOK, "exited")
}
func (h *Handler) UpdateInvoiceOrderItem(c echo.Context) error {
	type Req struct {
		Serial        int
		Qnt           float64
		Price         float64
		QntAntherUnit float64
		ItemName      string
	}

	req := new(Req)
	if err := c.Bind(req); err != nil {
		return err
	}
	rows, err := h.db.Raw("EXEC StkTrInvoiceDetailsUpdate  @Serial = ? , @Price = ? , @Qnt = ? ,  @QntAntherUnit = ? ,  @ItemName = ? ", req.Serial, req.Price, req.Qnt, req.QntAntherUnit, req.ItemName).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	return c.JSON(http.StatusOK, "updated")
}
func (h *Handler) CloseInvoiceOrder(c echo.Context) error {
	type Req struct {
		Serial int
		DocNo  int
	}
	code := c.Get("empCode")
	req := new(Req)
	if err := c.Bind(req); err != nil {
		return err
	}

	_, err := h.db.Raw("EXEC StkTrInvoiceHeadUpdate  @Serial = ? ,@UpdateTotalCash = ? , @AuditCode = ?", req.Serial, 1, code).Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "closed")
}
