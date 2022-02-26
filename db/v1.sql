
ALTER TABLE StkTrInvoiceHead
ADD DeletedAt DATETIME DEFAULT NULL; 

ALTER TABLE StkTrInvoiceHead
ADD Deleted BIT DEFAULT 0; 

GO
CREATE PROC StkTrInvoiceHeadDelete(@Serial INT)
AS 
BEGIN
    UPDATE StkTrInvoiceHead SET DeletedAt = GETDATE() , Deleted = 1 WHERE "Serial" = @Serial
END

GO
ALTER PROCEDURE [dbo].[StkTrInvoiceHeadList](@EmpCode int , @Finished bit = 0 , @Deleted bit = 0  ,@DateFrom  VARCHAR(20) = '',@DateTo  VARCHAR(20) = '' )
AS
BEGIN
DECLARE @SecLevel int
SET @SecLevel = (SELECT SecLevel FROM Employee WHERE EmpCode = @EmpCode)

SELECT  h.Serial ,  DocNo , DocDate ,h.EmpCode,
ISNULL (SUM(Qnt * Price), 0 ) TotalCash ,creator.EmpName ,AccountName,AccountCode,AccountSerial, Reserved ,  ISNULL(h.IsPosted,0) Finished
from StkTrInvoiceHead h
inner join 
StkTrInvoiceDetails d on h.Serial = d.HeadSerial
		JOIN Employee creator
		ON h.EmpCode = creator.EmpCode 
		join AccMs01 
		on h.AccountSerial = AccMs01.Serial 
		where 
		h.EmpCode = case when @seclevel < 4 then @Empcode else h.EmpCode end 
        AND h.DocDate <= case when @DateTo = '' THEN  h.DocDate ELSE  CONVERT(DATETIME, @DateTo, 102) END
        AND h.DocDate >= case when @DateFrom = '' THEN  h.DocDate ELSE  CONVERT(DATETIME, @DateFrom, 102) END
        AND h.Deleted = case when @Deleted = NULL THEN  h.Deleted ELSE @Deleted  END
		AND h.IsPosted = @Finished
        
group by HeadSerial ,DocNo,h.EmpCode,AccountSerial, h.Serial,DocDate,EmpName,AccountCode,AccountName, Reserved,IsPosted
END


