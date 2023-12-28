package summary

/*
SELECT walletAddress, kolAddress, mode, SUM(SellValue - Margin) FROM f_future_trading
WHERE CONVERT_TZ(timeStamp, '+00:00', '+08:00') >= '2023-12-01 00:00:00'
  AND CONVERT_TZ(timeStamp, '+00:00', '+08:00') < '2023-12-02 00:00:00'
  AND orderType in ("MARKET_CLOSE_FEE", "SL_ORDER_FEE", "TP_ORDER_FEE", "MARKET_LIQUIDATION")
GROUP BY walletAddress, kolAddress, mode;
*/
