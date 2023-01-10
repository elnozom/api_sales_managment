package handler

import (
	"net/http"
	"sms/model"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) ConvertToEta(c echo.Context) error {

	var resp int
	serial, _ := strconv.ParseInt(c.Param("serial"), 0, 64)

	err := h.db.Raw("EXEC StkTr01ConvertInvoice @Serial = ? ", serial).Row().Scan(&resp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) OrdersList(c echo.Context) error {

	var resp []model.OrderResp

	rows, err := h.db.Raw("EXEC StkTr01List").Rows()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var rec model.OrderResp
		err = rows.Scan(&rec.Serial, &rec.DocNo, &rec.DocDate, &rec.Discount, &rec.TotalCash, &rec.TotalTax)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		resp = append(resp, rec)
	}
	return c.JSON(http.StatusOK, resp)
}
