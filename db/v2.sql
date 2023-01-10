USE [MCSTREE]
GO

CREATE PROCEDURE StkTr01List(@StoreCode int )
	AS
    BEGIN
        SELECT 
        o.DocNo , o.DocDate , o.TotalCash ,o.StoreCode , s.StoreName 
        FROM StkTr01 o 
        JOIN StoreCode s ON o.StoreCode = s.StoreCode 
        WHERE o.StoreCode = CASE WHEN @StoreCode = 0 THEN o.StoreCode ELSE @StoreCode END
        AND isnull(o.TotalCash ,0) > 0
    END



GO
CREATE PROCEDURE StkTr01ConvertInvoice(@HeadSerial int )
	AS
    BEGIN
        DECLARE @DocNo nvarchar(16) 
        DECLARE @DocDate datetime 
        DECLARE @TransSerial int 
        DECLARE @StoreCode smallint 
        DECLARE @AccountSerial int
        DECLARE @TotalCash real 
        DECLARE @OrderNo int  
        DECLARE @IsPrinted bit  
        DECLARE @IsPosted bit 
        DECLARE @StkTr01Serial int 
        DECLARE @delivery_Date datetime 
        DECLARE @BuyOrSale bit 
        DECLARE @ShipCode int 
        DECLARE @CurCode int 
        DECLARE @CurRatio money
        DECLARE @TotalByCur money 
        DECLARE @EmpCode int 
        DECLARE @AuditCode int 
        DECLARE @Vat money


        SELECT  @DocNo = DocNo ,
                @DocDate = DocDate , 
                @TransSerial = TransSerial , 
                @StoreCode = StoreCode , 
                @AccountSerial = AccountSerial , 
                @TotalCash = TotalCash , 
                @OrderNo = OrderNo , 
                @IsPrinted = IsPrinted , 
                @IsPosted = IsPosted , 
                @StkTr01Serial = HeadSerial , 
                @delivery_Date = delivery_Date , 
                @BuyOrSale = DocNo , 
                @ShipCode = DocNo , 
                @CurCode = CurCode , 
                @CurRatio = CurRatio , 
                @TotalByCur = DocNo , 
                @EmpCode = DocNo , 
                @AuditCode = DocNo , 
                @Vat = DocNo 
                   FROM StkTr01 WHERE Serial = @HeadSerial

        INSERT INTO   StkTrInvoiceHead
        (DocNo, DocDate, TransSerial, StoreCode, AccountSerial, TotalCash, OrderNo,EmpCode,AuditCode ,IsPrinted, IsPosted, StkTr01Serial, delivery_Date, BuyOrSale, ShipCode, 
                                CurCode, CurRatio, TotalByCur , Vat)
        VALUES (@DocNo ,@DocDate , @TransSerial , @StoreCode , @AccountSerial , @TotalCash , @OrderNo , @EmpCode,@AuditCode ,@IsPrinted, @IsPosted, @StkTr01Serial, @delivery_Date, @BuyOrSale, @ShipCode, @CurCode, @CurRatio, @TotalByCur , @Vat )
    END

