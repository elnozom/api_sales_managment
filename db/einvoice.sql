USE MCSTREE;

DROP TABLE StkTrEInvoiceHead;
CREATE TABLE StkTrEInvoiceHead(
	"Serial" int IDENTITY(1,1) NOT NULL,
	internalID VARCHAR(100) NOT NULL,
	storeCode int NOT NULL,
	dateTimeIssued datetime NOT NULL  DEFAULT GETDATE(),
	totalDiscountAmount REAL NOT NULL DEFAULT 0,
	totalAmount REAL NOT NULL,
	totalTax REAL NOT NULL,
	stkTr01Serial int NULL,
	accountSerial int NOT NULL,
	deletedAt datetime NULL,
	deleted bit NULL
)
DROP TABLE StkTrEInvoiceDetails;
CREATE TABLE StkTrEInvoiceDetails(
	"Serial" int IDENTITY(1,1) NOT NULL,
	itemType VARCHAR(10) NOT NULL DEFAULT('EGS'),
	itemCode VARCHAR(100) NOT NULL,
	unitType VARCHAR(10) NOT NULL DEFAULT('BOX'),
	unitValue REAL NOT NULL,
	quantity INT NOT NULL,
	totalTaxableFees REAL NOT NULL,
	itemsDiscount REAL NOT NULL DEFAULT(0),
	itemSerial int NOT NULL,
	deletedAt datetime NULL,
	deleted bit NULL
)




GO
ALTER PROCEDURE StkTr01List
	AS
    BEGIN
        SELECT  o.Serial ,  o.DocNo ,  o.DocDate ,  ISNULL(o.Discount , 0) Discount  ,  o.TotalCash ,  ISNULL(o.TotalTax , 0) TotalTax   FROM StkTr01 o WHERE o.TotalCash IS NOT NULL AND ISNULL(EtaConverted , 0) = 0 AND TransSerial = 30
    END


GO
ALTER PROCEDURE StkTr01ConvertInvoice(@Serial int )
	AS
    BEGIN
        DECLARE @internalID VARCHAR(100)
        DECLARE @storeCode int
        DECLARE @dateTimeIssued datetime
        DECLARE @totalDiscountAmount REAL
        DECLARE @totalAmount REAL
        DECLARE @totalTax REAL
        DECLARE @accountSerial int
		
        SELECT  @internalID =   DocNo ,
                @storeCode = StoreCode , 
                @dateTimeIssued = DocDate , 
                @totalDiscountAmount = Discount , 
                @totalAmount = TotalCash , 
                @totalTax = ISNULL(TotalTax , 0)  , 
                @accountSerial = AccountSerial FROM StkTr01 WHERE Serial = @Serial
        INSERT INTO StkTrEInvoiceHead (internalID ,storeCode ,dateTimeIssued ,totalDiscountAmount ,totalAmount ,totalTax ,stkTr01Serial ,accountSerial) VALUES (@internalID ,@storeCode ,@dateTimeIssued ,@totalDiscountAmount ,@totalAmount ,@totalTax ,@Serial ,@accountSerial)



        INSERT INTO StkTrEInvoiceDetails ( itemCode  , unitValue , quantity , totalTaxableFees , itemsDiscount , itemSerial )
        SELECT  BarCodeUsed , Price , Qnt , Tax , Discount , ItemSerial 
        FROM    StkTr02 WHERE HeadSerial = @Serial  


		UPDATE StkTr01 SET EtaConverted = 1 WHERE "Serial" = @Serial
		SELECT SCOPE_IDENTITY() id
    END
